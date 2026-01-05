package main

import (
	"fmt"
	"log"

	"github.com/six-ddc/v2ex-go/internal/api"
)

func main() {
	client := api.NewClient()

	fmt.Println("Testing V2EX API...")
	fmt.Println("===================")

	// 测试获取主题列表
	fmt.Println("\n1. Testing GetTopicsByTab('all')...")
	topics, nodes, user, err := client.GetTopicsByTab("all")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Topics: %d\n", len(topics))
		fmt.Printf("   Nodes: %d\n", len(nodes))
		fmt.Printf("   User: %+v\n", user)

		if len(topics) > 0 {
			fmt.Println("\n   First 3 topics:")
			for i, t := range topics {
				if i >= 3 {
					break
				}
				fmt.Printf("   %d. %s\n", i+1, t.Title)
				fmt.Printf("      Node: %s, Author: %s, Replies: %d\n", t.Node.Name, t.Author, t.ReplyCount)
				fmt.Printf("      URL: %s\n", t.URL)
			}
		}

		if len(nodes) > 0 {
			fmt.Println("\n   First 5 nodes:")
			for i, n := range nodes {
				if i >= 5 {
					break
				}
				fmt.Printf("   - %s (%s)\n", n.Name, n.Code)
			}
		}
	}

	// 测试获取帖子详情
	if len(topics) > 0 {
		fmt.Println("\n2. Testing GetTopicDetail...")
		topic, replies, err := client.GetTopicDetail(topics[0].URL)
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			fmt.Printf("   Title: %s\n", topic.Title)
			fmt.Printf("   Author: %s\n", topic.Author)
			fmt.Printf("   Replies: %d\n", len(replies))
			fmt.Printf("   Content preview: %.100s...\n", topic.Content)

			if len(replies) > 0 {
				fmt.Println("\n   First 2 replies:")
				for i, r := range replies {
					if i >= 2 {
						break
					}
					fmt.Printf("   #%d @%s: %.50s...\n", r.Floor, r.Author, r.Content)
				}
			}
		}
	}

	fmt.Println("\n===================")
	fmt.Println("Test complete!")
}
