# Knative Quickstart CloudEvents Example

Simply put, the official Knative quickstart tutorial has quite a few gaps explaining things and it took me a while to piece things together. How can your app reply events to the broker?

This example deploys two services (```cloudevents-player``` and our custom app ```kn-app```) with two triggers. If you send an CloudEvent with type ```my-event``` from ```cloudevents-player```, ```kn-app``` will reply an CloudEvent with type ```app-event``` which will in turn be sent back to ```cloudevents-player```.

![drawio](https://user-images.githubusercontent.com/44191076/193742467-3f6af810-c47d-4360-93c7-eae01391a4b9.png)

The app is written in Golang using CloudEvents SDK, modified from the [second official serving example](https://github.com/knative/docs/tree/main/code-samples/serving/cloudevents/cloudevents-go). Both incoming and outputing events has a JSON payload with a single ```message``` field:

```json
{
  "message": "something"
}
```

## Run This Example

### Create Knative Quickstart Environment

Install Docker and download binaries of [kn](https://github.com/knative/client/releases), [kn-quickstart](https://github.com/knative-sandbox/kn-plugin-quickstart/releases) and [kind](https://github.com/kubernetes-sigs/kind/releases).

Start the demo environment:

```bash
./kn quickstart kind
```

### Generate and Upload kn-app Image

```bash
docker login
docker build . -t <your docker hub name>/kn-app -f Dockerfile
docker push <your docker hub name>/kn-app
```

Then modify the image source in ```kn-app.yaml``` so that it can find your image:

```yaml
containers:
- image: <your docker hub name>/kn-app:latest
```

### Deploy Resources in Knative

```bash
kubectl apply -f kn-app.yaml
```

### Wait for Services Ready

```
./kn service list
```

### Send and Receive Events

Open the URL of CloudEvents player (for example, ```http://cloudevents-player.default.127.0.0.1.sslip.io```) and send an event with type ```my-event```. The message is a JSON object with one field ```message```. You should see a second event with type ```app-event``` appear shortly.

<img width="958" src="https://user-images.githubusercontent.com/44191076/193742070-01103e74-f0f3-4b20-b16e-c711dbdf080d.png">

The event with ID ```42``` is the one replied by the app. You can modify the eveent ID, source and message in the code.

### ```K_SINK``` vs. Direct Reply

In ```kn-app.yaml``` the ```my-app``` has a environment variable ```K_SINK```, which will tell the code where is the event target. Here we set it with the URL of example-broker, which is the built-in broker in the Knative quickstart environment:

```
http://broker-ingress.knative-eventing.svc.cluster.local/default/example-broker
```

If ```K_SINK``` is not set (target is empty) and you send an direct [CloudEvent-compatible HTTP request](https://cloud.google.com/eventarc/docs/cloudevents) to the app, the code will reply an event directly as a HTTP response. This also works if you run the app as a local Docker container.
