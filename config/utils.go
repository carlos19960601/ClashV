package config

import (
	"fmt"
	"strings"

	"github.com/carlos19960601/ClashV/adapter/outboundgroup"
	"github.com/carlos19960601/ClashV/common/structure"
)

func trimArr(arr []string) (r []string) {
	for _, e := range arr {
		r = append(r, strings.Trim(e, " "))
	}
	return
}

func proxyGroupDagSort(groupsConfig []map[string]any) error {
	type graphNode struct {
		indegree  int
		outdegree int
		topo      int
		data      map[string]any
		option    *outboundgroup.GroupCommonOption
		from      []string
	}

	decoder := structure.NewDecoder(structure.Option{TagName: "group", WeaklyTypedInput: true})
	graph := make(map[string]*graphNode)

	for _, mapping := range groupsConfig {
		option := &outboundgroup.GroupCommonOption{}
		if err := decoder.Decode(mapping, option); err != nil {
			return fmt.Errorf("[代理组] %s: %s", option.Name, err.Error())
		}

		groupName := option.Name
		if node, ok := graph[groupName]; ok {
			if node.data != nil {
				return fmt.Errorf("[代理组] %s 名称重复", groupName)
			}

			node.data = mapping
			node.option = option
		} else {
			graph[groupName] = &graphNode{
				indegree:  0,
				outdegree: 0,
				topo:      -1,
				data:      mapping,
				option:    option,
				from:      nil,
			}
		}

		for _, proxy := range option.Proxies {
			if node, ex := graph[proxy]; ex {
				node.indegree++
			} else {
				graph[proxy] = &graphNode{
					indegree:  1,
					outdegree: 0,
					topo:      -1,
					data:      nil,
					option:    nil,
					from:      nil,
				}
			}
		}
	}

	// 拓扑排序
	queue := make([]string, 0)
	for name, node := range graph {
		if node.indegree == 0 {
			queue = append(queue, name)
		}
	}

	for ; len(queue) > 0; queue = queue[1:] {
		name := queue[0]
		node := graph[name]
		if node.option != nil {
			if len(node.option.Proxies) == 0 {
				delete(graph, name)
				continue
			}

			for _, proxy := range node.option.Proxies {
				child := graph[proxy]
				child.indegree--
				if child.indegree == 0 {
					queue = append(queue, proxy)
				}
			}
		}

		delete(graph, name)
	}

	// 没有循环引用的情况
	if len(graph) == 0 {
		return nil
	}

	// 有循环，定位循环
	for _, node := range graph {
		if node.option == nil {
			continue
		}

		if len(node.option.Proxies) == 0 {
			continue
		}

	}

	return nil

}
