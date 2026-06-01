# %% Imports
import asyncio
from ddgs import DDGS
from crawl4ai import AsyncWebCrawler, CrawlerRunConfig
from crawl4ai.markdown_generation_strategy import DefaultMarkdownGenerator

# %% DuckDuckGo Search Tool
# Find out if there's a chunk splitter designed for crawl4ai results
# Also... there's probably a better way to do this function in Go
def duckduckgo_search(prompt: str, max_results) -> list:
    result = DDGS().text("python programming", max_results=max_results)
    urls = [r["href"] for r in result]

    return urls

# %% Crawler
async def main():
    config = CrawlerRunConfig(
        markdown_generator=DefaultMarkdownGenerator()
    )
    async with AsyncWebCrawler() as crawler:
        result = await crawler.arun("https://github.com/crdcamp", config=config)

        if result.success:
            print("Raw Markdown Output:\n")
            print(result.markdown)  # The unfiltered markdown from the page
        else:
            print("Crawl failed:", result.error_message)

if __name__ == "__main__":
    asyncio.run(main())
