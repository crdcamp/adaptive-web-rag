# %% Imports
from ddgs import DDGS

import os
import sys
import psutil
import asyncio
import requests

from typing import List
from crawl4ai import AsyncWebCrawler, BrowserConfig, CrawlerRunConfig, CacheMode
from crawl4ai import DefaultMarkdownGenerator
from crawl4ai.content_filter_strategy import PruningContentFilter

import json

__location__ = os.path.dirname(os.path.abspath(__file__))
__output__ = os.path.join(__location__, "output")

# Append parent directory to system path
parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.append(parent_dir)

# Define save location for results
crawl_data_path = "crawl_data"
os.makedirs(crawl_data_path, exist_ok=True)

# Read prompt output from `llm-server.go`
with open (f"{crawl_data_path}/user_prompt.md") as f:
    prompt = f.read()

# %% DuckDuckGo Search Tool
async def duckduckgo_search(search_query: str, max_results: int) -> dict:
    print(f"Gathering {max_results} URLs for DuckDuckGo search: {search_query}")
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, lambda: DDGS().text(search_query, max_results=max_results))

    # Only extract the hrefs (urls)
    hrefs = [r["href"] for r in result]
    print(f"DuckDuckGo search complete for search: {search_query}")
    return hrefs

# Will need to eventually log the search URLs with a time stamp
# Need to add credit for this function
# %% Crawler
async def crawl_parallel(urls: List[str], max_concurrent: int = 4) -> dict:
    # === CLEAN SEMANTIC CONFIG (no links + pruning) ===
    prune_filter = PruningContentFilter(threshold=0.26, min_word_threshold=10)

    md_generator = DefaultMarkdownGenerator(
        content_filter=prune_filter,
        options={"ignore_links": True, "ignore_images": True}
    )

    crawl_config = CrawlerRunConfig(
        cache_mode=CacheMode.BYPASS,
        markdown_generator=md_generator,
        excluded_tags=["nav", "footer", "header", "script", "style"],
        word_count_threshold=50,
    )

    results_data = {}
    print("\n=== Parallel Crawling with Browser Reuse + Memory Check ===")

    # We'll keep track of peak memory usage across all tasks
    peak_memory = 0
    process = psutil.Process(os.getpid())

    def log_memory(prefix: str = ""):
        nonlocal peak_memory
        current_mem = process.memory_info().rss  # in bytes
        if current_mem > peak_memory:
            peak_memory = current_mem
        print(f"{prefix} Current Memory: {current_mem // (1024 * 1024)} MB, Peak: {peak_memory // (1024 * 1024)} MB")

    # Minimal browser config
    browser_config = BrowserConfig(
        headless=True,
        verbose=False,
        extra_args=["--disable-gpu", "--disable-dev-shm-usage", "--no-sandbox"],
    )

    # Create the crawler instance
    crawler = AsyncWebCrawler(config=browser_config)
    await crawler.start()

    try:
        # We'll chunk the URLs in batches of 'max_concurrent'
        success_count = 0
        fail_count = 0
        for i in range(0, len(urls), max_concurrent):
            batch = urls[i : i + max_concurrent]
            tasks = []

            for j, url in enumerate(batch):
                # Unique session_id per concurrent sub-task
                session_id = f"parallel_session_{i + j}"
                task = crawler.arun(url=url, config=crawl_config, session_id=session_id)
                tasks.append(task)

            # Check memory usage prior to launching tasks
            log_memory(prefix=f"Before batch {i//max_concurrent + 1}: ")

            # Gather results
            results = await asyncio.gather(*tasks, return_exceptions=True)

            # Check memory usage after tasks complete
            log_memory(prefix=f"After batch {i//max_concurrent + 1}: ")

            # Evaluate and return results as a dictionary
            for url, result in zip(batch, results):
                if isinstance(result, Exception):
                    print(f"Error crawling {url}: {result}")
                    fail_count += 1
                elif result.success:
                    results_data[url] = {
                        "href": url,
                        "content": getattr(result.markdown, 'fit_markdown', result.markdown)
                    }
                    success_count += 1
                else:
                    fail_count += 1

        results_data = list(results_data.values())

        print(f"\nSummary:")
        print(f"  - Successfully crawled: {success_count}")
        print(f"  - Failed: {fail_count}")

    finally:
        print("\nClosing crawler...")
        await crawler.close()
        # Final memory log
        log_memory(prefix="Final: ")
        print(f"\nPeak memory usage (MB): {peak_memory // (1024 * 1024)}")

    return results_data

async def main():
    urls = await duckduckgo_search(prompt, 10)
    if urls:
        print(f"Found {len(urls)} URLs to crawl")
        result = await crawl_parallel(urls, max_concurrent=10)

        # Save as a file meant to be overwritten.
        with open(f"{crawl_data_path}/crawl_results.json", "w") as f:
            json.dump(result, f)
        print(f"Web crawl results saved to: {crawl_data_path}/crawl_results.json")
    else:
        print("No URLs found to crawl")

if __name__ == "__main__":
    asyncio.run(main())
