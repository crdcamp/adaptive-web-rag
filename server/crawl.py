# %% Imports
import asyncio
from ddgs import DDGS
from crawl4ai import AsyncWebCrawler

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
async def main(search_query, max_results: int):
    hrefs = await duckduckgo_search(search_query, max_results)
    async with AsyncWebCrawler() as crawler:
        result = await crawler.arun_many(urls=hrefs)
        return [r.markdown.fit_markdown for r in result]

if __name__ == "__main__":
    result = asyncio.run(main("benefits and drawbacks of llama.cpp library", 8))
    print(result)
