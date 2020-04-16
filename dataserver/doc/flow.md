##  流程定义

使用Json格式定义流程，这个流程不同于审批流程，是用来定义内部的服务执行的过程的，通过流程定义各个服务的调用关系，在之前的年代称为服务编排，但这里还是使用流程这个不太准确的定义。

流程定义暂时不去考虑调用外部服务或者外部模块的情况，先考虑调用本系统内部注册的各个数据服务的情况，其实都是Restful服务应该处理方式是类似的，这里没有使用golang的rpc服务，也是考虑和通用Restful相兼容的问题。

BPMN2是个好东西，当初使用Activiti的时候其实是在使用这个规范，他是面向XML的，忒麻烦了。在golang世界里面有一些使用json定义流程的项目，比如：https://github.com/go-workflow/go-workflow  看了觉得很不错，但我不需要这么复杂，我只需要能够按照预定义的逻辑去调用我的服务然后返回结果即可，所以还是自己定一些数据结构来支持。

一个流程定义包含一个全局唯一名称“name”属性，必须包含一个start节点，流程处理程序会从start节点开始运行流程。

### start节点

定义流程的入口节点，包含启动流程时需要定义的变量和参数。

* 变量：就是在节点内定义的变量只能由本节点和其子节点访问
* 参数：启动流程时需要传入的参数值，所有节点都可以访问
* 数据类型：为了简化，我们只处理
  * string  字符串类型
  * number  数值类型，管他是整型还是浮点型，都用它
  * datetime 时间类型
  * array 数组
  * map  golang中的字典，js中的对象

``` json
{
  "name": "流程名称",
  "start": {
    "variables": {
      "参数1": {
        "type": "string|number|datetime|array|map",
        "value": {}
      },
      "参数1": {
        "type": "string|number|datetime|array|map",
        "value": {}
      }
    },
    "flow": [
      {
        "gate": "to",
        "activity1": {
          "style": "InnerService|Message|Loop",
          "params": {},
          "flow": []
        }
      }
    ]
  }
}
```

  ### flow 节点

start节点可以包含一个flow子节点，不能多了，只能有一个，如果没有也不行，因为那样这个流程就相当于什么都没做了，flow节点定义流程该往哪去。

flow节点是一个对象或者是一个数组：

* 对象：定义流程下一步flow。

  ```json
  {
  	"gate": "to",
      "activity1": {
      	"style": "InnerService|Message|Loop|Stdout",
          "params": {},
          "flow": [
              {
                  "gate":"ifto",
                  "if":"表达式",
                  "then":{},
                  "else":{}
              },
              {
                  "gate":"ifto",
                  "if":"表达式",
                  "then":{},
                  "else":{}
              },
              {
                  "gate":"end"             
              },
              {
                  "gate":"break"
              },
              {
                  "gate":"loop",                
                  "assign":["",""],
                  "while":"",
                  "step":["",""],
                  "do":{
                      "activity2":{}
                  }
              }
          ]
      }
  }
  ```

  ​	

* 数组：可以定义一堆下一步的控制逻辑，系统会按照他们出现的顺序一个一个执行。存在种情况就是前面执行的逻辑已经将流程终止或者返回的情况，那么这些后面的控制将不会被执行。

  #### gate=to

  gate属性定义当前flow的控制类型，包含to、ifto；gate节点外的其他属性定义了活动，属性名即为活动名称。控制类型为to表明flow将直接运行下面的活动，在to下面可以定义多个活动，flow将按照他们出现的顺序一个一个执行。如果某一个活动终止了当前流程，那么他后面的活动也不会被运行。

  #### gate=ifto

  条件判断，根据if属性定义的表达式的值确定flow执行的节点，如果表达式为true执行then属性定义的活动，如果表达式值为false执行else属性定义的活动

  #### gate=end

  结束处理，表示当前活动执行结束，交还控制权，流程会执行flow的下一个节点，或直接结束

  #### gate=break

  终止流程执行，立即退出

  #### gate=loop

  循环执行活动

  

#### 活动

活动activity类似BPMN中的概念，就是一个处理实际业务的单元，当然也包含一些用于控制流程状态的活动。`

```json
"activity1": {
    	"style": "InnerService|Message",
        "params": {},
        "flow": {
        	"ifto": {
            	"if": "Expression",
              	"then": {},
              	"else": {}
            },
            "to": {}
     	}
    }
```



属性名就是活动的名称，属性值定义活动的行为，一个活动至少包含style属性和flow属性，style定义活动的类型，flow定义活动执行结束后流程去向何处。params属性是可选的，定义活动需要用到的参数。

#### style=InnerService

调用系统内部定义的服务，就是在G_Service表中保存的那些服务。G_Service服务调用需要用到三组参数，:action、QueryString、PostBody

```json
"params":{
    "surl":"服务地址namespace.context的形式",
    "action":"请求的操作",
    "querystring":{
        "_repstyle":"map"
    }
    "body":{}
}
```

在params中可以使用采用${变量名}的方式来引用当前流程中的变量。