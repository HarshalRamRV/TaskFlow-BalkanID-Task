import PropTypes from "prop-types";

const TaskList = ({
  tasks,
  onDragStart,
  updateTaskStatus,
  deleteTask,
  status,
}) => {
  const onDragOver = (e) => {
    e.preventDefault();
  };

  const onDrop = (e, newStatus) => {
    const taskId = e.dataTransfer.getData("taskId");
    const taskTitle = e.dataTransfer.getData("taskTitle");
    const role = e.dataTransfer.getData("role");
    const group = e.dataTransfer.getData("group");
    updateTaskStatus(taskId, taskTitle, newStatus, role, group);
  };

  return (
    <div
      className="mx-2 w-full min-h-[10rem] bg-slate-200 p-5"
      onDragOver={(e) => onDragOver(e)}
      onDrop={(e) => onDrop(e, status)}
    >
      <h2 className="text-lg font-semibold">{status} Tasks</h2>
      <ul>
        {tasks
          .filter((task) => task.status === status)
          .map((task) => (
            <li
              key={task.id}
              className="flex flex-col border-b p-2 bg-white "
              draggable
              onDragStart={(e) => onDragStart(e, task.id, task.title, task.role, task.group)}
            >
              <div className="flex justify-between">
                <span
                  className="cursor-pointer"
                >
                  {task.title}
                </span>
                <button
                  onClick={() => deleteTask(task.id, task.role, task.group)}
                  className="text-red-500 hover:text-red-600 ml-2"
                >
                  Delete
                </button>
              </div>
              <div className="flex">
                <span>Role: {task.role}</span>
                <span className="px-2">Group: {task.group}</span>
              </div>
            </li>
          ))}
      </ul>
    </div>
  );
};

TaskList.propTypes = {
  tasks: PropTypes.array.isRequired,
  onDragStart: PropTypes.func.isRequired,
  updateTaskStatus: PropTypes.func.isRequired,
  deleteTask: PropTypes.func.isRequired,
  status: PropTypes.string.isRequired,
};

export default TaskList;
