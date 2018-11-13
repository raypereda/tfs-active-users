# tfs-active-users
Team Foundation Server (TFS) utility that lists users that have made at least 10 change sets in the last year.

I use this to count the number of active user at Molina Healthcare. [Semmle](https://semmle.com/), the security source code company, only charges for active users!

Sample output:
```
c:\ tfs-active-users.exe

executing the command line below
tf history /collection:http://tfs.molina.mhc:8080/tfs/HSS/ $/HSS /recursive /version:D2017-11-12~D2018-11-12 /noprompt
344 active users with at least 10 check-ins

Check-ins      User
--------- -----------------
   21149    Siruvuru, Murali
   13144    TFSSERVICE
   2993    SQLDLSvc
   1230    Gnanadhiraviam...
   861    Kumar, Anand
   691    Kumar, Pradeepa
   476    Lee, Hyuk
   393    Pawar, Ganeshrao
   388    Kandagaddala, ...
```
