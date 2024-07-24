# kubecnf

Add and remove kubeconfigs from the main config file

## Usage

**Info:** with unspecified config file (-c / --config) default value is used

### Add a new cluster config
~~~
kubecnf -c /path/to/main/config add -f /path/to/new/cluster/config
~~~
### Remove a cluster config
~~~
kubecnf -c /path/to/main/config remove cluster-name
~~~
### Rollback to the previous config
~~~
kubecnf -c /path/to/main/config rollback
~~~