# radiko-linker

radiko-linker is a service that generates links to programs distributed by radiko.  
A link to the program is generated by specifying the station, broadcast date, and time.  

# usage

The user accesses "/make-link" and enters the area code, day of the week, and broadcast time to get the URL to access the radio program.

![make-link](img/make-link.drawio.svg)

The obtained URL is in the Web API format of make-linker, so you can open radiko by accessing the URL.

![jump](img/jump.drawio.svg)
