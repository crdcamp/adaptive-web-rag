# %% Imports
from ddgs import DDGS

import os
import sys
import psutil
import asyncio
import requests

from typing import List
from crawl4ai import AsyncWebCrawler, BrowserConfig, CrawlerRunConfig, CacheMode

__location__ = os.path.dirname(os.path.abspath(__file__))
__output__ = os.path.join(__location__, "output")

# Append parent directory to system path
parent_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.append(parent_dir)

# NEED TO ADD TIMER FOR BOTH OF THESE FUNCTIONS
# In a later implementation, you should also consider using sessions
# that run in tandem with llama-server

# %% DuckDuckGo Search Tool
# Find out if there's a chunk splitter designed for crawl4ai results
# Also... there's probably a better way to do this function in Go (are we gonna rewrite this in go too?)
async def duckduckgo_search(search_query: str, max_results: int) -> list:
    print(f"Gathering {max_results} URLs for DuckDuckGo search: {search_query}")
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, lambda: DDGS().text(search_query, max_results=max_results))
    hrefs = [r["href"] for r in result]
    print(f"DuckDuckGo search complete for search: {search_query}")
    return hrefs

# %% Crawler
# Credit: https://github.com/coleam00/ottomator-agents/blob/main/crawl4AI-agent/crawl4AI-examples/3-crawl_docs_FAST.py
async def crawl_parallel(urls: List[str], max_concurrent: int = 3):
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
        verbose=False,   # corrected from 'verbos=False'
        extra_args=["--disable-gpu", "--disable-dev-shm-usage", "--no-sandbox"],
    )
    crawl_config = CrawlerRunConfig(cache_mode=CacheMode.BYPASS)

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

            # Evaluate and return results
            results_data = []
            for url, result in zip(batch, results):
                if isinstance(result, Exception):
                    print(f"Error crawling {url}: {result}")
                    fail_count += 1
                elif result.success:
                    results_data.append(result)
                    success_count += 1
                else:
                    fail_count += 1

                return results_data

        print(f"\nSummary:")
        print(f"  - Successfully crawled: {success_count}")
        print(f"  - Failed: {fail_count}")


    finally:
        print("\nClosing crawler...")
        await crawler.close()
        # Final memory log
        log_memory(prefix="Final: ")
        print(f"\nPeak memory usage (MB): {peak_memory // (1024 * 1024)}")


async def main():
    urls = await duckduckgo_search("benefits and drawbacks of llama.cpp library", 8)
    if urls:
        print(f"Found {len(urls)} URLs to crawl")
        results = await crawl_parallel(urls, max_concurrent=10)
        print(type(results))
    else:
        print("No URLs found to crawl")

if __name__ == "__main__":
    asyncio.run(main())
