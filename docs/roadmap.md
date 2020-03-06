## Road map 

### Pick project name:
* ChurnDr
* PodX-ray
* PodChurn (current name of CRD)
* 

#### Principle
* Keep it simple in it's base form to encourage adoption. Advanced topology can be configured when needed.  
* currently only a controller and a frontend UI deployment are needed.


### v0.5
* filter by labels.
* Topic grouping.
* UI scrolling.
* Alerting on pods waiting to be scheduled for extended periods of time. 
* Purge old pod history on max-pods history reached 
  Keeping Last 15 restarts with 4Kb log limit and a global pod limit of 250,000 pods. Memory usage is 16Gigabytes.
  Using 250k pod limit on a cluster with 100 namespaces. It should be able to keep history 2,500 different pods per namespace.


### v1

* Flow diagram
* Create docker CI and Continuos Delivery to Artifactory
* Create IKS deployment Yaml.
* Test in Dog foot and FDP pre-prod clusters
* Finish defining CRD (Custom Resource Definition)
* Option to record logs from all/selection of containers/sidecars in pod not only failed one.
* Persistance:
    - Fully in memory with memory usage limits. (limit implemented in max history per-pod, Need limit in max number of pods).
    - Persist periodically to disk and read data on pod restart if no other persistance is defined. 
    - *Optional* Persist driver to redis and/or MySQL.
* (Optional) in CRD* Save pod termination logs to specified S3 bucket.
* (Optional) API throttle limit.
* Improve Unit-testing and code coverage.

#### v1.1
* Bug fixes
* Look for more users besides FDP.

### v1.2 OSS
* CI/CD in travisCI with dockerhub publishing.
* boot-strap,demo and presentation.

### v2

* Add support for other types of notifications (Detect errors on Ingress?, Storage, Disk full mounts? )
* Alerting on deployments unable to grow to capacity for extended time.
* Frontend integration with DevPortal or IKSM. We still want a stand-alone UI to keep it modular testable and retain open source path.
* Notification integration with other enterprise tools and devportal.
* Look at improving noise reduction and grouping of errors.
* UI review 
* Auth  UI and API



NOTES:
- Redis 
  Standalone : Short outage with no data loss.
  Elastic Cache : Managed no data loss. 
  Redis Kubernetes Operator ( https://product.spotahome.com/redis-operator-for-kubernetes-released-9c6775e7da9b).

