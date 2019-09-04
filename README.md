# shotify
A screenshot service written in Go using chromedp package.

The design of the service is modular in sense I have a screenshot service listening on one port and accepts POSTs with either a list
of urls or a file with new tab seperated urls, which will be processed by service and shares the results as a json. This response, if 
successful fetch happened, shares the link of the image file to be sent as  GET request on an archiving server. The archiver/image
server will handle all requests for images.

Since the data is not sensitive, I have avoided adding authentication check for requests which basically elads to anyone having
the response json to be able to fetch the screenshot. But as a small security mechanism, I have shared only link of the image and not
the hostname:port of image server. The idea is that the clients of the service will have RESTful endpoints along with both servers'
ip and port to be able to issue requests.

As a practical way for my testing and sharing, I have intantiated 2 servers from teh same program one listening at 8002 port and another 
at 8008 for screenshot request and for image request respectively.



How to run:
 clone the repository.
 go get the required packages
 run go build main.go 
 install chrome headless:
          sudo apt-get update
          sudo apt-get install -y libappindicator1 fonts-liberation
          wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
          sudo dpkg -i google-chrome*.deb
 ./main
 
 
Basically workflow is as below:
  1. A screenshot service, shotify, is listening on <ip>:8002 for incoming requests as below format:
          
            * curl -X POST 'http://127.0.0.1:8002/?urls=<comma seperated url list>'     
              eg/- curl -X POST 'http://127.0.0.1:8002/?urls=netflix.com, http://google.com'
              
            *curl -H "Content-Type: text/plain" -X POST 'http://127.0.0.1:8002/?urlfile' --data-binary @<Filename>
              eg/- curl -H "Content-Type: text/plain" -X POST 'http://127.0.0.1:8002/?urlfile' --data-binary @testurls.txt
              
  2. As a response, in case of success we get below json:
  
          [
          
              {
              
               "url":"https://www.google.com",
              
              "status":"success",
              
              "link":"https%3A%2F%2Fwww.google.com"
              
              },
              
              {
              "url":"http://www.github.com",
              
              "status":"success",
              
              "link":"http%3A%2F%2Fwww.github.com"
              
              }
              ]
              
   3. Using the link in the json, one can issue a GET request for file as,    
   
            * curl http://127.0.0.1:8008/?fileName=<valid_link>
               eg/- curl http://127.0.0.1:8008/?fileName=http%3A%2F%2Fwww.github.com > somefile.png
            

 Scalability for upto 1 million screenshots a day:
  I have used local file system on the archive server to save screenshot images. Database is not very well suited to store files and 
  since all we need is a POST and GET on the image files alone, we get better performance using filesystem.
  But caveat was some file systems support a fixed number of subdirectories(like In ext3 a directory can have at most
  32,000 subdirectories). To handle this, I have used a directory tree structure like root/<outerfolder>/<inner_folder>filename where
  outer and inner are calculated using hash of the filename. Scales pretty well. But this again hits inherent limit of disk capacity.
  To handle it, we must run archive server behind a proxy/ load balancer so with some consistent hashing so that we can distribute the load
  across multiple machines . This combined with directory tree can handle any number of load with horizontal scaling.
  
  
Limitations on my implementation:
    1. The lookup cache used to avoid repeated fetches is crude and needs to be improved to prevent lock starvations
    2. Hardcoded extensions(.png), ip addresses of archive server. Should be made configurable using a json file.\
    3. Hacky and crude url parse logic using url.Parse(). since chromedp rejected urls without http/s, added that too in hardcoded way
    4. To handle installations of chrome headless , should have made a docker compose file 
    5. Vendoring support using go modules to be added.
    6. Load testing is not done.
 
 
