# 如何使用本工具

## 介绍

  这个工具用于清洗 [Ownthink](https://www.ownthink.com/) 的知识图谱 RDF 数据，将它变成属性图模型。产出结果为一个 vertex.csv 文件和 edge.csv 文件, 分别是清洗后的顶点数据和边数据，用于neo4j数据库。目前只对数据进行了简单去重及去除顶点中的换行回车。
  
  该工具参考并改写于
  https://github.com/jievince/rdf-converter
  原工具支持 Nebula ，不能直接用于neo4j

## 如何使用

使用 --path 参数指定知识图谱的三元组数据的路径

go build

生成rdf2neo4j或rdf2neo4j.exe

rdf2neo4j --path ownthink_v2.csv

生成
vertex.csv 存放顶点

edge.csv   存放边

之后可以使用 neo4jadmin import进行导入了
