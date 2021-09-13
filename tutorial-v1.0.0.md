
## API 1
Responsibility: Apply pipeline

|Name | Details                                        |                  
|---|-----------------------------------------------|
|Id |1                                             | 
|API Version |  v1 | 
|Url | [http://host:port/api/v1/pipelines?url=[repo]&revision=[commitId/branch]&purging=[ENABLE/DISABLE]]()       |
|Request Type |  POST |                            
|Tekton Version |  v1alpha1 |

#### Payload

##### Basic build and push
```

{
  "name": "test",
  "steps": [
    {
      "name": "build",
      "type": "BUILD",
      "service_account": "sa_name",
      "input": {
        "type": "git"
      },
      "outputs": [{
        "type":"image",
        "url":"image_url",
        "revision":"revision"
      }
      ]
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
      "name": "build",
      "type": "BUILD",
      "service_account": "sa_name",
      "input": {
        "type": "git"
      },
      "outputs": [{
        "type":"image",
        "url":"image_url",
        "revision":"revision"
      }
      ],
      "arg":{
        "data":{
          "arg1":"value1"
        }
      }
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
      "name": "build",
      "type": "BUILD",
      "service_account": "sa_name",
      "input": {
        "type": "git"
      },
      "outputs": [{
        "type":"image",
        "url":"image_url1",
        "revision":"revision"
      },
        {
          "type":"image",
          "url":"image_url2",
          "revision":"revision"
        }
      ],
      "arg":{

        "configMaps":[
          {
            "name":"configmap_name",
            "namespace":"namespace_name"
          }
        ]
      }
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
      "name": "build",
      "type": "BUILD",
      "service_account": "sa_name",
      "input": {
        "type": "git"
      },
      "outputs": [{
        "type":"image",
        "url":"image_url1",
        "revision":"revision"
      },
        {
          "type":"image",
          "url":"image_url2",
          "revision":"revision"
        }
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
      "name": "build",
      "type": "BUILD",
      "service_account": "test-sa",
      "input": {
        "type": "git"
      },
      "outputs": [
        {
          "type": "image",
          "url": "zeromsi2/spring-boot-image:1.0.0",
          "revision": "1.0.0"
        }
      ],
      "arg": {
        "configMaps": [
          {
            "name": "cm-test",
            "namespace": "tekton"
          }
        ]
      }
    },
    {
      "name": "deploy",
      "type": "DEPLOY",
      "outputs": [
        {
          "deployment_resource": {
            "agent": "local_agent",
            "name": "ubuntu",
            "namespace": "default",
            "type": "deployment",
            "replica": 2,
            "images": [
              {
                "image_index": 0,
                "image": "zeromsi2/spring-boot-image:1.0.0"
              },
              {
                "image_index": 1,
                "image": "ubuntu"
              }
            ]
          }
        }
      ]
    }
  ]
}

```


---
**NOTE**
Inorder to register any new agent add agent info for key ```AGENTS``` in environment like , 
[name]&&[url]&&[token] separated by comma (,)
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

