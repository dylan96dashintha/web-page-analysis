# web-page-analysis
This project allows you to analyze web pages through a simple UI. Uses a Go server and a static frontend (index.html)

## Getting Started
### 1.Clone the repo
```bash
git clone [<your-github-repo-url>](https://github.com/dylan96dashintha/web-page-analysis.git)
```

### 2. Navigate to the Project Directory
```bash
cd [<repo>](https://github.com/dylan96dashintha/web-page-analysis.git)
```

### 3. Build the Docker image
```bash
docker build -t analysis-img -f Dockerfile .
```

### 4. Run the docker container
```bash
docker run -d -p 8080:8080 --name web-analysis-container analysis-img
```

### 5. Open the Frontend
open the index.html (static/index.html) file in your browser 

