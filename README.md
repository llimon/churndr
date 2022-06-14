## ChurnDoctor (ChurnDr)

## Problem Statement

Kubernetes does a excelent job at keeping applications running, some time even when they are having serious stability issues. Kubernetes Will stop sending traffic to pods failing to respond readiness probes and terminate pods failing liveness probes. 

This creates several secondary problems.

1. It masks and hides service affecting bugs. That need to be fixed if they are not dealt with then can accumulate to a unsustainable level.
2. Because Kubernetes manages pod lifecycle is difficult to know if your application is serving traffic at proper capacity or is experiencing issues not yet reported by users.
3. When kubernetes Restarts pods it marks as misbehaving logs are overwritten by new instance of pod scheduled. 


## Solutions

* X-ray back in time of pod crashes and restarts
* Alerting, reporting and realtime monitoring of pod issues 
* Alert Grouping and Noise reduction 
* 