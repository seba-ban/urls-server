# URLS

Finding myself constantly flooded with loads of open tabs on my macBook, I created this very simple server to manage them. It can be used together with automator workflow to save all open tabs at once. The project is not polished, I've just quickly written it just so it works more or less...

## Local Deployment

Build the image:

```bash
docker image build -t urls-server --target production .
```

Two most important things:

* container will listen on port 8080;
* server uses sqlite db located in /db/urls.db within the container – create a mount if needed.

So run the container (adjust as needed):

```bash
docker container run \
  -d --name urls-server \
  --restart unless-stopped \
  -p 28080:8080 \
  -v ~/Documents/dbs/urls.db:/db/urls.db \
  urls-server
```

Now you can create an Automator workflow consisting of two steps:

![automator workflow](img/automator.png)

* run javascript – the below code can be used and adjusted if needed:

  ```js
  function run(input, parameters) {
    const urls = [];
      for (const browserName of ["Safari", "Chrome"]) {
        try {
          const browser = Application(browserName);
          for (let i = 0; i < browser.windows.length; i++)
            for (let j = 0; j < browser.windows[i].tabs.length; j++) {
              const url = browser.windows[i].tabs[j].url() || "";
              const description = browser.windows[i].tabs[j].name() || "";
              if (url) urls.push({ url, description });
            }
        } catch(err) {
          console.log(err);
        }
      }
      return JSON.stringify({data: urls.sort()});
  }
  ```

  The above code loops over all tabs in all windows of Safari and Chrome, and gathers urls and tab names. The returned json will be passed as an input to the next step.

* run shell script with input passed to stdin:

  ```bash
  curl -XPOST -H'Content-Type:application/json' --data-binary @- --no-buffer http://localhost:28080/urls
  ```

  Basically, we're sending the json generated by the javascript in previous step to the urls server using curl.
