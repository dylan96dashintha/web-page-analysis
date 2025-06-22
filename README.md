# web-page-analysis
This project allows you to analyze web pages through a simple UI. Uses a Go server and a static frontend (index.html)

## Project Overview
The app analyses a webpage by
  * Identifying internal and external links.
  * Checking accessibility of each link (via HTTP status codes).
  * Skipping duplicate links to reduce redundant checks.
  * Usage of a worker pool to handle performance bottlenecks when checking the accessibility.

Note: JavaScript-rendered content (e.g., Instagram, Facebook) is not currently supported.

## Technologies Used
Backend -Go
Frontend - HTML, CSS
Devops -Docker

## Prerequisites
Docker
Go (only required for local dev, not for running via docker)
Git

## External Dependencies
github.com/PuerkitoBio/goquery
github.com/gorilla/mux
github.com/sirupsen/logrus
gopkg.in/yaml.v3
Standard Go libraries

These dependencies will be installed automatically via go mod tidy

## Getting Started

#### 1.Clone the repo
```bash
git clone https://github.com/dylan96dashintha/web-page-analysis.git
```

#### 2. Navigate to the Project Directory
```bash
cd web-page-analysis
```

#### 3. Build the Docker image
```bash
docker build -t analysis-img -f Dockerfile .
```

#### 4. Run the docker container
```bash
docker run -d -p 8080:8080 --name web-analysis-container analysis-img
```

#### 5. Open the Frontend
open the index.html (static/index.html) file in your browser manually.

## Main Assumptions

#### Internal/External Link Classification
If a link's hostname differs from the base URL, it is considered external, otherwise it's internal.

#### Accessibility Check
The system sends HTTP requests to each link. If no errors occur (timeouts) and the status code is within the 200–299 range, the link is considered accessible.

#### Duplicate Removal
Duplicate links are ignored to prevent redundant processing.

## Challenges & Solutions

High number of links increases API response latency - Implemented a worker pool to parallelize accessibility checks
Some URLs take too long to respond - Applied TCP timeouts for outbound HTTP requests
CORS errors when calling the backend from browser - Added CORS handling middleware to the Go server
Inconsistent href formatting	- Normalized and joined relative paths with the base URL

## Imporvements

Introduce pagination for the result - when there are multiple inaccessibility links
Add support for JavaScript-rendered content (Instagram, Facebook)

