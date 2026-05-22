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

# Run the server 

```terminal
chmod +x llama-server.sh
./llama-server.sh
```
