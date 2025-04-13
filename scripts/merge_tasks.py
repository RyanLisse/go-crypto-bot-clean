import json
import os

# Read the existing tasks
with open('tasks/tasks.json', 'r') as f:
    existing_tasks = json.load(f)

# Read the new tasks
with open('tasks/tasks.json.new', 'r') as f:
    new_tasks = json.load(f)

# Create a map of existing tasks by ID for quick lookup
existing_tasks_map = {task['id']: task for task in existing_tasks['tasks']}

# Merge the tasks, preserving status and subtasks of existing tasks
merged_tasks = {
    'tasks': [
        {
            **new_task,
            'status': existing_tasks_map.get(new_task['id'], {}).get('status', new_task['status']),
            'subtasks': existing_tasks_map.get(new_task['id'], {}).get('subtasks', new_task.get('subtasks', []))
        }
        for new_task in new_tasks['tasks']
    ]
}

# Write the merged tasks to a new file
with open('tasks/tasks.json.merged', 'w') as f:
    json.dump(merged_tasks, f, indent=2)

print('Tasks merged successfully!')
