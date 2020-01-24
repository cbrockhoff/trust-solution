package policyTrees

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mettegroenbech/policyGateway/serviceWhitelist"
)

type PolicyTrees map[string]*Node

type Node struct {
	name       string
	parameter  bool
	matchClaim bool
	allow      map[string]bool
	children   []*Node
}

func BuildPolicyTrees(whitelist serviceWhitelist.ServiceWhitelist) PolicyTrees {
	policyTrees := make(map[string]*Node)

	for _, entry := range whitelist.Entries {
		urlPath := strings.Split(strings.Trim(entry.URL, "/"), "/")
		allowMap := make(map[string]bool)
		for _, a := range entry.Allow {
			allowMap[a] = true
		}

		method := strings.ToUpper(entry.Method)
		_, found := policyTrees[method]
		if !found {
			root := Node{}
			policyTrees[method] = &root
		}

		longestChain, remainingUrlPath := findLongestChain(policyTrees[method], urlPath)

		tempPolicyMap := make(map[string]*Node)
		for i, p := range remainingUrlPath {
			node := &Node{}
			if longestChain.name == p && (!longestChain.parameter || (longestChain.parameter && isParameter(p))) {
				node = longestChain
			} else {
				node.name = p
			}

			if isParameter(p) {
				node.parameter = true
				if isClaimMatch(p) {
					node.matchClaim = true
				}
			}

			if i == len(remainingUrlPath)-1 {
				node.allow = allowMap
			}

			tempPolicyMap[p] = node
		}
		for i, p := range remainingUrlPath {
			if i < len(remainingUrlPath)-1 {
				nodeToGiveChild := tempPolicyMap[p]
				child := tempPolicyMap[remainingUrlPath[i+1]]

				nodeToGiveChild.children = append(nodeToGiveChild.children, child)
			}
		}
		if longestChain.name != remainingUrlPath[0] {
			longestChain.children = append(longestChain.children, tempPolicyMap[remainingUrlPath[0]])
		}
	}

	for _, n := range policyTrees {
		sortNodeChildren(n)
	}

	return policyTrees
}

func sortNodeChildren(node *Node) {
	if node.children == nil {
		return
	}
	sort.Slice(node.children, func(i, j int) bool {
		return !node.children[i].parameter && node.children[j].parameter
	})
	for _, c := range node.children {
		sortNodeChildren(c)
	}
}

func IsRequestAllowed(policyTrees map[string]*Node, method string, requestPathString string, caller string, claims map[string]string) bool {
	requestPath := strings.Split(requestPathString, "/")
	method = strings.ToUpper(method)
	node, found := policyTrees[method]
	if !found {
		return false
	}
	for i, p := range requestPath {
		if i == len(requestPath)-1 {
			_, found := node.allow[caller]
			return found
		}
		if (node.name == p && !node.parameter) || node.parameter { // check if it is a matchClaim
			if node.parameter && node.matchClaim {
				pureName := strings.TrimSuffix(strings.TrimPrefix(node.name, "{:"), "}")
				claim, found := claims[pureName]
				if !found || (found && claim != p) {
					return false
				}
			}
			if node.children == nil {
				return false
			} else {
				for _, child := range node.children {
					if (child.name == requestPath[i+1] && !child.parameter) || child.parameter {
						node = child
						break
					}
				}
			}
		} else {
			return false
		}
	}
	return false
}

func Print(policyTrees map[string]*Node) {
	for method, node := range policyTrees {
		fmt.Println(method)
		policyPrint(node, 1)
		fmt.Println()
	}
}

func policyPrint(node *Node, level int) {
	name := node.name
	if name == "" {
		name = "root"
	}
	leading := ""
	for i := 0; i < level; i++ {
		leading = leading + "-"
	}
	fmt.Println(leading, name, "parameter:", node.parameter, "matchClaim:", node.matchClaim, "allow:", node.allow)
	if node.children != nil {
		for _, child := range node.children {
			policyPrint(child, level+3)
		}
	}

}

func findLongestChain(root *Node, urlPath []string) (*Node, []string) {
	node := &Node{}
	foundRootChild := false
	if root.children != nil {
		for _, child := range root.children {
			if (child.name == urlPath[0] && !child.parameter) || (child.parameter && isParameter(urlPath[0]) && child.matchClaim == isClaimMatch(urlPath[0])) { // should also check whether the urlPath is parameter and whether it is claimMatch
				node = child
				foundRootChild = true
				break
			}
		}
		if !foundRootChild {
			return root, urlPath
		}
	} else {
		return root, urlPath
	}
	for i, p := range urlPath {
		if i == len(urlPath)-1 {
			return node, urlPath[i:]
		}
		if (node.name == p && !node.parameter) || (node.parameter && isParameter(p) && node.matchClaim == isClaimMatch(p)) { // should also check whether the urlPath is parameter and whether it is claimMatch
			if node.children == nil {
				return node, urlPath[i:]
			} else {
				for _, child := range node.children {
					if (child.name == urlPath[i+1] && !child.parameter) || (child.parameter && isParameter(urlPath[i+1]) && child.matchClaim == isClaimMatch(urlPath[i+1])) { // should check whether it is claimMatch
						node = child
						break
					}
				}
			}
		} else {
			return node, urlPath[i:]
		}
	}
	return root, urlPath
}

func isParameter(p string) bool {
	return strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}")
}

func isClaimMatch(p string) bool {
	return strings.HasPrefix(p, "{:")
}
