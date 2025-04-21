'use client'
import { CheckIcon, XIcon } from "lucide-react";
import { useEffect, useState } from "react"
type Task = {
  id: string;
  userId: string;
  name: string;
  status: string;
  description: string;
  timeEstimate: number; // in hours
  dueDate: string; // in YYYY-MM-DD format
  isComplete: boolean;
}

export default function Home() {

  const [tasks, setTasks] = useState<Task[]>([])
  const [view, setView] = useState<'view' | 'create' | 'edit'>("view")
  const [taskToEdit, setTaskToEdit] = useState<Task>()

  useEffect(() => {
    try {
      fetch("http://localhost:8080/tasks")
        .then((res) => {
          if (!res.ok) throw new Error("Failed to fetch");
          return res.json();
        })
        .then((data) => setTasks(data))
        .catch((err) => console.error(err));
    } catch (error) {
      console.error("An unexpected error occurred:", error);
    }
  }, [])

  const exportCSVC = () => {
    fetch("http://localhost:8080/tasks/export", {
      method: "GET",
    })
      .then((res) => {
        if (!res.ok) {
          throw new Error("Failed to export tasks");
        }
        return res.blob(); // Receive the response as a Blob
      })
      .then((blob) => {
        const link = document.createElement("a");
        const url = URL.createObjectURL(blob);
        link.href = url;
        link.download = "tasks.csv"; // Set the filename for the download
        link.click();
        URL.revokeObjectURL(url); // Clean up the object URL
      })
      .catch((err) => {
        console.error("Error exporting tasks:", err);
      });
  }

  const handleCreateTask = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const newTask = {
      id: "", // Backend will assign ID
      name: formData.get("name"),
      status: "Pending",
      description: formData.get("description"),
    };

    fetch("http://localhost:8080/tasks", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(newTask),
    })
      .then((res) => {
        if (!res.ok) {
          throw new Error("Failed to create task");
        }
        return res.json();
      })
      .then((createdTask) => {
        setTasks((prev) => [...prev, createdTask]);
        setView("view");
      })
      .catch((err) => {
        console.error("Error creating task:", err);
      });
  }

  return (
    <div className="background-wave h-screen bg-opacity-50 bg-base-300">
      <div>
        <div className="container mx-auto p-4">
          <div className="flex  justify-between items-center mb-6 mt-8">
            <h1 className="text-5xl font-bold">
              go<span className="text-primary">Task</span>.
            </h1>
            <div className="flex justify-end gap-3 ">
              <button className="btn btn-lg btn-primary btn-soft" onClick={() => setView("create")}>
                + New Task
              </button>
              <button
                onClick={exportCSVC}
                className="btn btn-accent btn-lg  btn-soft"
              >
                Export Tasks as CSV
              </button>
            </div>
          </div>
          {view === "view" && (
            <div className="bg-base-100">

              <table className="table">
                <thead>
                  <tr >
                    <th>ID</th>
                    <th>Name</th>
                    <th>Status</th>
                    <th>Description</th>
                    <th>Time Estimate</th>
                    <th>Due Date</th>
                    <th>Is Complete</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody className="rounded-3xl">
                  {tasks.map((task: Task) => (
                    <tr key={task.id}>
                      <td>{task.id}</td>
                      <td>{task.name}</td>
                      <td>{task.status}</td>
                      <td>{task.description}</td>
                      <td>{task.timeEstimate || "N/A"}</td>
                      <td>{task.dueDate || "N/A"}</td>
                      <td>{task.isComplete ? "Yes" : "No"}</td>
                      <td className="flex  items-center gap-2">
                        <button
                          className="btn btn-error btn-sm  btn-soft"
                          onClick={() => setTasks(tasks.filter(t => t.id !== task.id))}
                        >
                          Delete
                        </button>
                        <button
                          className="btn btn-warning btn-sm  btn-soft"
                          onClick={() => {
                            setTaskToEdit(task);
                            setView("edit");
                          }}
                        >
                          Edit Task
                        </button>
                        <button
                          className={`${task.isComplete ? "btn btn-soft btn-sm" : "btn btn-sm  btn-soft"
                            } `}
                          onClick={() => {
                            fetch(`http://localhost:8080/tasks/complete/${task.id}`, {
                              method: "PUT",
                              headers: {
                                "Content-Type": "application/json",
                              },
                            })
                              .then((res) => {
                                if (!res.ok) {
                                  throw new Error("Failed to update task");
                                }
                                return res.json();
                              })
                              .then((savedTask) => {
                                setTasks((prev) =>
                                  prev.map((t) => (t.id === savedTask.id ? savedTask : t))
                                );
                              })
                              .catch((err) => {
                                console.error("Error updating task:", err);
                              });
                          }}
                        >
                          {task.isComplete ? <XIcon className="text-red-400" /> : <CheckIcon className="text-green-400" />}
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          {view === 'create' && (
            <div className="flex flex-col card bg-base-100 p-6 rounded-lg shadow-md">
              <form
                onSubmit={handleCreateTask}
              >
                <div>
                  <label htmlFor="name">Name:</label>
                  <input
                    type="text"
                    id="name"
                    name="name"
                    required
                    className="w-full input"
                  />
                </div>
                <div className="flex flex-col gap-2">
                  <div>
                    <label htmlFor="description">Description:</label>
                    <textarea
                      id="description"
                      name="description"
                      required
                      className="w-full textarea"
                    ></textarea>
                  </div>
                  <div>
                    <label htmlFor="timeEstimate">Time Estimate:</label>
                    <input
                      type="text"
                      id="timeEstimate"
                      name="timeEstimate"
                      className="w-full input"
                    />
                  </div>
                  <div>
                    <label htmlFor="dueDate">Due Date:</label>
                    <input
                      type="date"
                      id="dueDate"
                      name="dueDate"
                      className="w-full input"
                    />
                  </div>
                </div>
                <div className="flex justify-end mt-4">
                  <button
                    type="button"
                    className="bg-gray-500 text-white px-4 py-2 rounded mr-2"
                    onClick={() => setView("view")}
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="bg-blue-500 text-white px-4 py-2 rounded"
                  >
                    Save
                  </button>
                </div>
              </form>
            </div>
          )}
          {view === "edit" && taskToEdit && (
            <div className="flex flex-col card bg-base-100 p-6 rounded-lg shadow-md">
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  const formData = new FormData(e.currentTarget);

                  const updatedTask = {
                    ...taskToEdit,
                    name: taskToEdit.name,
                    description: taskToEdit.description,
                    timeEstimate: taskToEdit.timeEstimate,
                    dueDate: taskToEdit.dueDate,
                  };

                  fetch(`http://localhost:8080/tasks/${taskToEdit.id}`, {
                    method: "PUT",
                    headers: {
                      "Content-Type": "application/json",
                    },
                    body: JSON.stringify(updatedTask),
                  })
                    .then((res) => {
                      if (!res.ok) {
                        throw new Error("Failed to update task");
                      }
                      return res.json();
                    })
                    .then((savedTask) => {
                      setTasks((prev) =>
                        prev.map((t) => (t.id === savedTask.id ? savedTask : t))
                      );
                      setView("view");
                    })
                    .catch((err) => {
                      console.error("Error updating task:", err);
                    });
                }}
              >
                <div className="mb-4">
                  <label
                    htmlFor="name"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
                  >
                    Name:
                  </label>
                  <input
                    type="text"
                    id="name"
                    name="name"
                    required
                    value={taskToEdit.name}
                    onChange={(e) => {
                      setTaskToEdit((prev) => ({
                        ...prev,
                        name: e.target.value,
                      }));
                    }}
                    className="input w-full"
                  />
                </div>

                <div className="mb-4">
                  <label
                    htmlFor="description"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
                  >
                    Description:
                  </label>
                  <textarea
                    id="description"
                    name="description"
                    required
                    value={taskToEdit.description}
                    onChange={(e) => {
                      setTaskToEdit((prev) => ({
                        ...prev,
                        description: e.target.value,
                      }));
                    }}
                    className="textarea w-full"
                  />
                </div>

                <div className="mb-4">
                  <label
                    htmlFor="timeEstimate"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
                  >
                    Time Estimate:
                  </label>
                  <input
                    type="text"
                    id="timeEstimate"
                    name="timeEstimate"
                    value={taskToEdit.timeEstimate}
                    onChange={(e) => {
                      setTaskToEdit((prev) => ({
                        ...prev,
                        timeEstimate: e.target.value,
                      }));
                    }}
                    className="input w-full"
                  />
                </div>

                <div className="mb-4">
                  <label
                    htmlFor="dueDate"
                    className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
                  >
                    Due Date:
                  </label>
                  <input
                    type="date"
                    id="dueDate"
                    name="dueDate"
                    value={taskToEdit.dueDate}
                    onChange={(e) => {
                      setTaskToEdit((prev) => ({
                        ...prev,
                        dueDate: e.target.value,
                      }));
                    }}
                    className="input w-full"
                  />
                </div>

                <div className="flex justify-end mt-6 gap-4">
                  <button
                    type="button"
                    className="btn btn-error btn-soft"
                    onClick={() => setView("view")}
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    className="btn btn-success btn-soft "
                  >
                    Save
                  </button>
                </div>
              </form>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

