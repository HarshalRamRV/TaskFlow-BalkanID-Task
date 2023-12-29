import { useState, useEffect, useContext } from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";
import TaskList from "./TaskList";
import Navbar from "./NavBar";
import { AuthContext } from "../auth/AuthContext";
import CsvUpload from './CsvUpload';

function TaskBoard() {
  const [tasks, setTasks] = useState([]);
  const [newTask, setNewTask] = useState("");
  const [role, setRole] = useState("");
  const [group, setGroup] = useState("");
  const [user, setUser] = useState(null);
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");


  const navigate = useNavigate();
  const { token } = useContext(AuthContext);

  if (!token) {
    navigate("/");
  }

  useEffect(() => {
    fetchTasks();
  }, []);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await axios.get("http://localhost:3001/user", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.status === 200) {
          setUser(response.data);
        }
      } catch (error) {
        console.error(error);
      }
    };

    fetchUser();
  }, [token]);
  const fetchTasks = () => {
    axios.get("http://localhost:3001/tasks").then((response) => {
      setTasks(response.data);
      console.log(response.data);
    });
  };

  const addTask = () => {
    axios
      .post("http://localhost:3001/tasks", {
        title: newTask,
        status: "To-Do",
        role: role,
        group: group,
      })
      .then(() => {
        setNewTask("");
        fetchTasks();
      });
  };

  const updateTaskStatus = (id, title, newStatus, role, group) => {

    if (tasks.length === 0) {
      console.log("Tasks are not loaded yet.");
      return;
    }

  
    if (user.role === role && user.group === group) {
      axios
        .put(`http://localhost:3001/tasks/${id}`, {
          title: title,
          status: newStatus,
        })
        .then(() => {
          fetchTasks();
        })
        .catch((error) => {
          console.error("Failed to update task:", error);
        });
    } else {
      alert("Unauthorized to update task.");
    }
  };
  
  const deleteTask = (id,role,group) => {
  
    if (user.role === role && user.group === group) {
      axios
        .delete(`http://localhost:3001/tasks/${id}`)
        .then(() => {
          fetchTasks();
        })
        .catch((error) => {
          console.error("Failed to delete task:", error);
        });
    } else {
      alert("Unauthorized to delete task.");
    }
  };


  const onDragStart = (e, taskId, taskTitle, role, group) => {
    e.dataTransfer.setData("taskId", taskId);
    e.dataTransfer.setData("taskTitle", taskTitle);
    e.dataTransfer.setData("role", role);
    e.dataTransfer.setData("group", group);
  };

  return (
    <>
      <Navbar user={user} />
      <div className="container w-[90vw] mt-4 m-auto">
        <div className="flex space-x-2 flex justify-center ">
          <input
            type="text"
            placeholder="New Task"
            className="border p-2 flex-grow"
            value={newTask}
            onChange={(e) => setNewTask(e.target.value)}
          />
          <button
          onClick={addTask}
          className="bg-blue-500 hover:bg-blue-600 text-white p-2 rounded"
        >
          Add Task
        </button>
        </div>
        <div className="flex space-x-2 mt-4 flex justify-center">
          <input
            type="text"
            placeholder="Role"
            className="border p-2 flex-grow"
            value={role}
            onChange={(e) => setRole(e.target.value)}
          />
          <input
            type="text"
            placeholder="Group"
            className="border p-2 flex-grow"
            value={group}
            onChange={(e) => setGroup(e.target.value)}
          />
        </div>

        <div className="mt-4 flex w-full">
          <TaskList
            tasks={tasks}
            onDragStart={onDragStart}
            updateTaskStatus={updateTaskStatus}
            deleteTask={deleteTask}
            status="To-Do"
          />
          <TaskList
            tasks={tasks}
            onDragStart={onDragStart}
            updateTaskStatus={updateTaskStatus}
            deleteTask={deleteTask}
            status="Pending"
          />
          <TaskList
            tasks={tasks}
            onDragStart={onDragStart}
            updateTaskStatus={updateTaskStatus}
            deleteTask={deleteTask}
            status="Completed"
          />
        </div>
        <div className="flex flex-col	flex justify-center">
        <CsvUpload
          token={token}
          endpoint="http://localhost:3001/upload/users"
          onSuccess={(message) => setSuccessMessage(message)}
          onError={(message) => setErrorMessage(message)}
          uploadtype="Users"
        />
      
        <CsvUpload
          token={token}
          endpoint="http://localhost:3001/upload/tasks"
          onSuccess={(message) => setSuccessMessage(message)}
          onError={(message) => setErrorMessage(message)}
          uploadtype="Tasks"
        />
      
        {successMessage && <div className="text-green-500">{successMessage}</div>}
        {errorMessage && <div className="text-red-500">{errorMessage}</div>}
      </div>
      

      </div>
    </>
  );
}

export default TaskBoard;
