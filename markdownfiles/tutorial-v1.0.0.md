## API 1

Responsibility: Apply pipeline

|Name | Details                                        |                  
|---|-----------------------------------------------|
|Id |1                                             | 
|API Version |  v1 | 
|Url | [http://host:port/api/v1/pipelines?url=[repo]&revision=[commitId/branch]&purging=[ENABLE/DISABLE]]()       |
|Request Type |  POST |    
|Content Type |  JSON/XML |  
|Tekton Version |  v1alpha1 |

#### Payload

##### Basic build and push

```
{
  "name": "test",
  "steps": [
    {
      "name": "name",
      "type": "BUILD",
      "trigger": "AUTO",
      "params": {
        "repository_type": "git",
        "revision": "revision",
        "service_account": "service_account_name",
        "images": "image_url"
      },
      "next": []
    }
  ]
}

```

##### Build passing arguments and push

```
{
  "name": "test",
  "steps": [
    {
      "name": "name",
      "type": "BUILD",
      "trigger": "AUTO",
      "params": {
        "repository_type": "git",
        "revision": "revision",
        "service_account": "service_account_name",
        "images": "image_url1,image_url2",
        "arg": "key1:value1,key2:value2"
      },
      "next": [
      ]
    }
  ]
}
```

##### Build passing arguments using configmap and push

```
{
  "name": "test",
  "steps": [
    {
      "name": "name",
      "type": "BUILD",
      "trigger": "AUTO",
      "params": {
        "repository_type": "git",
        "revision": "revision",
        "service_account": "service_account_name",
        "images": "image_url1,image_url2",
        "args_from_configmaps": "namespace/config_map_name",
      },
      "next": [
        "deployDev"
      ]
    }
  ]
}

```

##### Build and push to multiple registry

````
{
  "name": "test",
  "steps": [
    {
      "name": "name",
      "type": "BUILD",
      "trigger": "AUTO",
      "params": {
        "repository_type": "git",
        "revision": "revision",
        "service_account": "service_account_name",
        "images": "image_url1,image_url2"
      },
      "next": [

      ]
    }
  ]
}

````

##### Build and deploy

```
{
  "name": "test",
  "steps": [
    {
      "name": "name",
      "type": "BUILD",
      "trigger": "AUTO",
      "params": {
        "repository_type": "git",
        "revision": "revision",
        "service_account": "service_account_name",
        "images": "image_url1"
      },
      "next": [
        "deployDev"
      ]
    },
    {
      "name": "deployDev",
      "type": "DEPLOY",
      "trigger": "AUTO",
      "descriptors":[],
      "params": {
        "agent": "agent_name",
        "name": "name of resource",
        "namespace": "k8s_resource_namespace",
        "type": "k8s_resource_name",
        "env": "env_name",
        "images": "image_url"
      },
      "next": null
    }
  ]
}

```

---
**NOTE**

- ``` trigger ``` can be ```AUTO``` / ```Manual```. If there are any k8s resources to be applied, add those as json or
  yaml in ```descriptors```

---

## API 2

Responsibility: Get pipeline logs

|Name | Details                                        |                  
|---|-----------------------------------------------|
|Id |2                                             | 
|API Version |  v1 | 
|Url | [http://host:port/api/v1/pipelines/[processId]?page=[page]&limit=[limit]]()       |
|Request Type |  GET |                            
|Tekton Version |  v1alpha1 |

## API 3

Responsibility: Get pipeline events

|Name | Details                                        |                  
|---|-----------------------------------------------|
|Id |2                                             | 
|API Version |  v1 | 
|Url | [ws://host:port/api/v1/pipelines/ws?[processId]]()       |
|Request Type |  GET |                            
|Tekton Version |  v1alpha1 |

