Here's the installation setup so far:

```bash
# Install llama.cpp
brew install llama.cpp

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
