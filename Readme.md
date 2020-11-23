# yougif

## Build the project

Make sure you have installed golang on your computer. Run
``bash
go build main.go
``

to build the project.

Execute the binary
``
./main
``

Please make sure that the folder "migrations", and the file "config.json" and "client_secret.json" are on the same level as the file "main" before you start the application.


## Setup Database

This is just an example and should not been used in production!
``
docker container run -p 5432:5432 -e POSTGRES_PASSWORD=password postgres
``

## Configuration

The project contains a file named "config.json". In this file are for example the database connection, or the number of videos to be downloaded.

### YouTube API

To retrieve the Most Popular Videos we need access to the YouTube API. To do so, we need to set up a Google project according to the following instructions.
"https://developers.google.com/youtube/v3/quickstart/go" 
Afterwards the file "client_sceret.json" must be put in the right place.

Start the application once and you will be prompted to open a link.
``
Go to the following link in your browser then type the authorization code:
https://accounts.google.com/o/oauth2/auth?access_type=offline...
``

Grant access to you application and paste the code from the browser into the console.

### server

``
server.port The port on which the rest API will be reachable
``

### database

``
database.host Hostname of the database
database.port Port of the database
database.user User of the database
database.password Password of the database user
``

### youTube

``
youTube.region Specifies the region from which the Top X videos should be downloaded
youTube.totalNumbersOfVideosToDownload How many videos should be downloaded
``

### processing

``
processing.numberOfWorkers Of many threads should process the videos into single images. Do not set this number too high
processing.ffmpegBinaryPath Where is the ffmpeg binary located. framesPerSecond /usr/local/ffmpeg
processing.framesPerSecond The number of images which should be created of one second from video (5 would result in one image every 0,2 seconds)
``

## API

This service provides an REST API

### search

Searches for the text "q" in the database and returns the FrameId, the YouTube VideoId and the timestamp. The FrameId can be used to query the image via the endpoint "/images/:id".

``
GET /search?q=hello
``

Response:
``json
{
    "data": [
        {
            "id": "4838b617-8bcd-484d-b5dc-e49299d4584a",
            "youTubeVideoId": "Rtmhv5qEBg0",
            "timeInMs": 313000
        }
    ]
}
``

### images

Returns the image as bytes (JPEG).

``
GET /images/:id
``
