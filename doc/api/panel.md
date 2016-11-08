1、GET  获取看板数据 /panel?favor=true&q=queryString
   response:{
     code:200,
     panels:[
       {
            id:int              -- 看板id    
            name:string         -- 看板名称，
            charts：int         -- 看板下图表数量
            update_at：time     -- 看板更新时间
        },
        ...
     ]  
  }

2、GET 获取看板包含的图表  /panel/chart/panel_id?q=queryString
   response:{
      code:200,
      data:[
       {
            chart_name:string,      — 图表名称
            chart_id: int           — 图表自增id
            update_at: time         — 图表更新时间
            refer_count: int        — 图表引用计数

       },
       ......
     ]
   }

3、GET 获取所有图表  /chart?page=2&pageSize=5&q=queryString
   response:{
      code:200,
      data:[
       {
          name:string           — 图表名称，
          id:int                — 图表id,
          chart_image:          — 图表展示图,
          update_at:DateTime    — 更新时间
          refer_count:int       — 引用计数
       },
       ......
      ]
   }
  
    page：        当前页码
    pageSize：    每页条数
    q:        查询字符串


———————————— 暂时不确定———————————
    
   GET  获取单个图表的数据和配置项信息---/chart/id
   response:{
      code:200,
      data:[
         {
            name:string        — 图表名称
        id:int        — 图表id
        size:int        — 图表尺寸
     }   
       ]
   }


——————————————————


4、PUT  创建新图表  /chart
   {
      name:string    — 图表名称
      size:enum        — 图表大小
      ele:[        — 图表元素
        {
            name:string        — 指标名称
            metric:string        - metrci
            tags: string        — tags
            thumbnail: string    — 缩略图数据，Base64字符串
        }
    ]
    }

5、POST 编辑已有图表 /chart
   {
      id: int        — 图表id
      name:string    — 图表名称
      size:enum        — 图表大小
      ele:[        — 图表元素
        {
            name:string        — 指标名称
            metric:string        - metrci
            tags: string        — tags
            thumbnail: string    — 缩略图数据，Base64字符串
        }
    ]
    }
    

6、POST 更新图表缩略图  /chart/thumbnail
{
    id:int             — chart id 
    thumbnail: string 
}


9、PUT 创建看板     /panel
    {
        name:"",--看板名称
        charts:[---新增图表id集
            {
                id:
                name:
                ....
            },
            ....
        ]
    }


11、POST 编辑看板 /panel
    {
        id:int
        name:string
        charts:[
            {
                id:
                name:
                ...
            }
       ]
     }

12、POST 删除在视窗展示的看板 /broad/favor/id--看板id
    response:{
      code:200
    }

13、DELETE 删除看板  /board/id
    
14、delete  删除图表与看板的关联  /relation/Pid（看板id）+/Cid(图表id)----仅删除图表与看板的关联关系，保留图表缩略图
   response:{
      code:200
   }
