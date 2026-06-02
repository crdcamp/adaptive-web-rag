# %% Imports
import asyncio
from ddgs import DDGS
from crawl4ai import AsyncWebCrawler

# NEED TO ADD TIMER FOR BOTH OF THESE FUNCTIONS

# %% DuckDuckGo Search Tool
# Find out if there's a chunk splitter designed for crawl4ai results
# Also... there's probably a better way to do this function in Go
async def duckduckgo_search(search_query: str, max_results: int) -> list:
    print(f"Gathering {max_results} URLs for DuckDuckGo search: {search_query}")
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, lambda: DDGS().text(search_query, max_results=max_results))
    hrefs = [r["href"] for r in result]
    print(f"DuckDuckGo search complete for search: {search_query}")
    return hrefs

# %% Crawler
async def main(search_query, max_results: int):
    hrefs = await duckduckgo_search(search_query, max_results)
    async with AsyncWebCrawler() as crawler:
        result = await crawler.arun_many(urls=hrefs)
        print(type(result))
    return result

if __name__ == "__main__":
    result = asyncio.run(main("benefits and drawbacks of llama.cpp library", 8))
    print(result[0])
