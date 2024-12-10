### 关于如何使用codesandbox
1. 下载docker相关环境: `https://www.docker.com/`

2. 配置好数据库相关信息

   1. 首先在创建好对应的库表。其中**有关题目的输入输出以及题目配置的表中必须包含以下字段**

      ```mysql
      id  bigint auto_increment comment 'id' primary key
      identity    varchar(36)   not null comment '唯一ID'
      judgeCase   text          null comment '判题用例（json 数组）'
      judgeConfig text          null comment '判题配置（json 对象）'
      created_at  datetime      not null comment '创建时间'
      updated_at  datetime      not null comment '更新时间'
      deleted_at  datetime      null comment '删除时间'
      ```

   2. 然后到`settings/application.dev.yaml`更改`MysqlConfig`的配置

      ```yaml
       MysqlDBName: "polaris-oj"
       MysqlDBPassword: "******"
      ```

      