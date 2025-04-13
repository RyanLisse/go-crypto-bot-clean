const fs = require('fs');

// Read the existing tasks
const existingTasks = JSON.parse(fs.readFileSync('tasks/tasks.json', 'utf8'));
// Read the new tasks
const newTasks = JSON.parse(fs.readFileSync('tasks/tasks.json.new', 'utf8'));

// Create a map of existing tasks by ID for quick lookup
const existingTasksMap = {};
existingTasks.tasks.forEach(task => {
  existingTasksMap[task.id] = task;
});

// Merge the tasks, preserving status and subtasks of existing tasks
const mergedTasks = {
  tasks: newTasks.tasks.map(newTask => {
    const existingTask = existingTasksMap[newTask.id];
    if (existingTask) {
      // Preserve status, subtasks, and other fields from existing task
      return {
        ...newTask,
        status: existingTask.status,
        subtasks: existingTask.subtasks || newTask.subtasks
      };
    }
    return newTask;
  })
};

// Write the merged tasks to a new file
fs.writeFileSync('tasks/tasks.json.merged', JSON.stringify(mergedTasks, null, 2));
console.log('Tasks merged successfully!');
