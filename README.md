# TO DO THIS IN PYTHON OR GO???

... It's still not too late to switch :)

But you could write just the crawl4ai part in python (obviously) and everything else in
Go... is that a bad idea?

Here's the installation setup so far:

```bash
# Setup go-llama.cpp
git clone --recurse-submodules https://github.com/go-skynet/go-llama.cpp
cd go-llama.cpp
make libbinding.a

# Setup Python environment and crawl4ai
python3 -m venv
source venv/bin/activate
pip install -r requirements.txt
crawl4ai-setup
```
