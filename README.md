# csv-webapp

![Build status][logo]

[logo]: docs/passing.svg

## What is csv-webapp?

csv-webapp was developed for the ThisSolution. Description of a web-application:

> Need an app that will ftp a daily file back with call CDR info, extract it and read/calculate costing based on the csv file.

> Each day the provider uploads a .tgz compressed file based on the date. We would like to pull it back, and then create a cleaned up xls file including some of the columns that are in the file, but also calculate a new ‘cost’ column based on the duration in one of the column. Another of the column specifies the tier that the pricing is charged using.

> I would like the app to be able to run in a docker container, can be php or another language.

> I would like to be able to open the web page (with a username / password), and the press a ‘Get Data button’. When this is pressed, the app will connect to the ftp server and pull back the ftp files that have not previously been pulled back (store the last file name previously pulled back), extract the tgz file and store all the raw data into a db.

> On a setting page, we should be able to specify a table that stores the field value for the different tiers, and then a $ rate for each tier. We also need to specify a minimum charge out $, a minimum duration, and the charge out rate per duration.

> I need to be able to press a generate report button, that will need a start and ending date. It will then need to pressent the data in a clean screen that will show the require columns (from the raw data), plus the calcualted ‘cost’ column. At the bottom we would like the total cost in a line  

## Development Status

csv-webapp is under agile development. There are still some parts that has to be done, but the main prototype is ready for testing.

**get report button and config page will be developed shortly**

## Architecture

![architecture][arch]

[arch]: docs/architecture.png

## Before using application

**NOTE: if you will be using docker-compose skip this part**
**NOTE: you have to download and install mongodb on your machine**

* After downloading and installing mongodb on your machine, create a database and collection using terminal. You need to do this step for authentication to work properly.

```
$ sudo systemctl start mongod

$ mongo

$ use cdr

$ db.createCollection("users")

$ db.users.insert({
    email: 'email@of.yours',
    password: 'your_password'
})
```

* Write down this email and password somewhere, you will need it to login to webapp

## Example Usage

* Download the repository using git commands:

```
$ git clone <repository>
```

* Download and configure [golang](https://golang.org/dl/).

* If you will be using Docker skip 4 step.

* On your machine run following commands:

```
$ go build -ldflags="-r -w"
  ./csv-webapp

//if you prefer scripts then run

$ chmod u+x build.sh
  ./build.sh
```

* For Docker build run following commands:

```
$ docker build -t csv-webapp_csv-webapp .

//then run

$ docker run -it -p 5050:5050 csv-webapp_csv-webapp

//or if you prefer container running in background

$ docker run -d -p 5050:5050 csv-webapp
```

* For Docker Compose run following command:

```
$ docker-compose up
```