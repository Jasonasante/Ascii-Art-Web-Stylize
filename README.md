# ascii-art-web

Description: This is basic web application made using golang and html that allows one to convert a string into a Graphical User Interface (GUI).

Authors: Abdul Raheem Khan , Jason Asante-Twumasi and Douglas Barco

Usage: how to run

1.Git clone the ascii-art-web repository

2.In your terminal, run "go run main.go" in the ascii-art-web repository

3.Now open a browser and enter "localhost:8080" into the URL

4.Then select your prefered banner, enter you text into the input box then press the convert button

5.A GUI of ur input text should appear in an output box.

How to run docker?

You have to build the program with the command :

docker build -t ascii-web-docker .

Then, run it with the command :

docker run -it --rm -p 8080:8080 ascii-web-docker

Implementation details: algorithm
