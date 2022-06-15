# TLS mailer

---

### this repository is for PRIVATE use of repository owner
##### no one else is allowed to use, copy or fork this repo for now.

---
The idea of this package comes from investigation of how 
smtp protocol works. It is build around SendMail function 
that uses implicit TLS encryption. Simple client is included,
which can be configured to use either STARTTLS or TLS.

### For those who are not familiar with details here is simplified explanation: 

There are two versions of smtp protocol: 
 - ESMTP with STARTTLS (commonly called STARTTLS), standart port 587
 - SMTPS (commonly called SSL), standart port 465

They differ by the way of establishing secure connection.
First one(STARTTLS) for long time was considered as 
recommended. During my investigation I've found that at 
last few years ago the situation finally has changed. I 
will not give links and RFCs here, you can find 
approvements by yourself.

Difference is that using ESMTP connection to smtp server is   
established at plaintext unencrypted mode. Then it 
negotiate to exchange acceptable cipher types and send 
STARTTLS command. After that if both sides have equal 
cipher types encryption starts. If not (plaintext 
connection can easily be attacked and ciphers can be 
modified by man-in-the-middle) - connection will be 
continued unencrypted. 

There are some "workarounds" to deal with this main problem
(recommendations for clients to inform about inability to 
use encryption, etc...) but they are just recommendations.
Implementation details depend on programmers of each client.
And that is only main problem, there are another, less 
dangerous, but anyway...

SMTPS that is commonly called SSL(illogical, it also uses
modern TLS encryption, just like STARTTLS and HTTPS for
example) encrypts connection from the beginning. If it was 
unable to set encrypted connection - data exchange does not
even start. 

Also SMTPS was claimed as deprecated almost just after it 
was accepted(in 1999). Some smtp servers accept SMTPS 
connections (Google Mail server for example), some do not 
(Outlook mail server). Just recently implicit TLS 
encryption became recommended by RFCs way of smtp protection.
---
"net/smtp" builtin package contains default SendMail function 
that internally uses STARTTLS mechanism.
In the honor of golang developers I should mention that 
default SendMail function from "net/smtp" package returns 
error if it was unable to encrypt after STARTTLS command.
---
As about me - personally I think that "full" encryption just
from the beginning is better. That was the reason to create
this package for my own email sending.

This package contains adopted for implicit TLS SendMail 
function - the core idea. It can be used through included 
client or directly. Client can be set to fallback to STARTTLS
if smtp server does not accept implicit TLS.

To make TLS SendMail function maximally interchangeable with 
"net/smtp"."SendMail" - signature is just the same.

Actually my TLS SendMail function was made by exact copying
of "net/smtp"."SendMail" function (at 01.06.2022) with change 
of tcp dial from "net"."Dial()" to "tls"."Dial()" and 
commented out starttls command sending (not necessary over 
already encrypted connection). Some other small changes were
made because of using unexported "net/smtp"."Client" fields 
and methods that were inside "net/smtp"."SendMail" function.