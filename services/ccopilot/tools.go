package ccopilot

type toolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function toolCallFunction `json:"function"`
}

type toolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type toolDefinition struct {
	Type     string      `json:"type"`
	Function functionDef `json:"function"`
}

type functionDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

func toolsDefinitions() []toolDefinition {
	return []toolDefinition{
		{
			Type: "function",
			Function: functionDef{
				Name:        "read_file",
				Description: "Read the content of a file from the repository's default branch",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "Path to the file relative to repository root (e.g. src/main.go)",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: functionDef{
				Name:        "list_directory",
				Description: "List files and directories in a repository path",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "Directory path relative to repository root (empty string for root)",
						},
					},
					"required": []string{"path"},
				},
			},
		},
		{
			Type: "function",
			Function: functionDef{
				Name:        "search_code",
				Description: "Search for code in the repository by pattern",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "Search term or regex pattern to find in the codebase",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: functionDef{
				Name:        "get_diff",
				Description: "Get the diff between two branches or commits",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"base": map[string]any{
							"type":        "string",
							"description": "Base branch name or commit SHA",
						},
						"head": map[string]any{
							"type":        "string",
							"description": "Head branch name or commit SHA",
						},
					},
					"required": []string{"base", "head"},
				},
			},
		},
		{
			Type: "function",
			Function: functionDef{
				Name:        "write_file",
				Description: "Create or update a file in a new branch. This will also create the branch if needed.",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "Path to the file relative to repository root (e.g. src/main.go)",
						},
						"content": map[string]any{
							"type":        "string",
							"description": "Full file content to write",
						},
						"branch": map[string]any{
							"type":        "string",
							"description": "Branch name to write to (will be created if it doesn't exist). Use a descriptive name starting with ccopilot-",
						},
						"message": map[string]any{
							"type":        "string",
							"description": "Commit message describing the change",
						},
					},
					"required": []string{"path", "content", "branch", "message"},
				},
			},
		},
		{
			Type: "function",
			Function: functionDef{
				Name:        "create_branch",
				Description: "Create a new branch from the default branch",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"name": map[string]any{
							"type":        "string",
							"description": "Name for the new branch. Use a descriptive name starting with ccopilot-",
						},
					},
					"required": []string{"name"},
				},
			},
		},
		{
			Type: "function",
			Function: functionDef{
				Name:        "create_pull_request",
				Description: "Create a pull request from a head branch to the default branch",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"title": map[string]any{
							"type":        "string",
							"description": "Title for the pull request",
						},
						"body": map[string]any{
							"type":        "string",
							"description": "Body/description for the pull request",
						},
						"head": map[string]any{
							"type":        "string",
							"description": "Head branch name (the branch with changes)",
						},
					},
					"required": []string{"title", "body", "head"},
				},
			},
		},
	}
}
