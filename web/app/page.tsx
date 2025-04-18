'use client'
import { useEffect, useState } from "react"
type Task = {
    id: string;
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

    return (
        <>
            {view === "view" && (
                <div>
                    <div className="flex justify-end">
                        <button className="bg-blue-500 text-white px-4 py-2 rounded" onClick={() => setView("create")}>
                            Create New
                        </button>
<button
    onClick={() => {
        // Trigger the CSV export by making a GET request to the backend
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
    }}
    className="bg-green-500 text-white px-4 py-2 rounded mt-4"
>
    Export Tasks as CSV
</button>

                    </div>
                    <table>
                        <thead>
                            <tr>
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
                        <tbody>
                            {tasks.map((task: Task) => (
                                <tr key={task.id}>
                                    <td>{task.id}</td>
                                    <td>{task.name}</td>
                                    <td>{task.status}</td>
                                    <td>{task.description}</td>
                                    <td>{task.timeEstimate || "N/A"}</td>
                                    <td>{task.dueDate || "N/A"}</td>
                                    <td>{task.isComplete ? "Yes" : "No"}</td>
                                    <td>
                                        <button
                                            className="bg-red-500 text-white px-2 py-1 rounded mr-2"
                                            onClick={() => setTasks(tasks.filter(t => t.id !== task.id))}
                                        >
                                            Delete
                                        </button>
                                        <button
                                            className="bg-yellow-500 text-white px-2 py-1 rounded mr-2"
                                            onClick={() => {
                                                setTaskToEdit(task);
                                                setView("edit");
                                            }}
                                        >
                                            Update
                                        </button>
                                        <button
                                            className={`${
                                                task.isComplete ? "bg-green-500" : "bg-gray-500"
                                            } text-white px-2 py-1 rounded`}
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
                                            {task.isComplete ? "Mark as Incomplete" : "Mark as Complete"}
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
            {view === 'create' && (
                <div>
                    <form
                        onSubmit={(e) => {
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
                          }}
                    >
                        <div>
                            <label htmlFor="name">Name:</label>
                            <input
                                type="text"
                                id="name"
                                name="name"
                                required
                                className="border px-2 py-1 rounded"
                            />
                        </div>
                        <div>
                            <label htmlFor="description">Description:</label>
                            <textarea
                                id="description"
                                name="description"
                                required
                                className="border px-2 py-1 rounded"
                            ></textarea>
                        </div>
                        <div>
                            <label htmlFor="timeEstimate">Time Estimate:</label>
                            <input
                                type="text"
                                id="timeEstimate"
                                name="timeEstimate"
                                className="border px-2 py-1 rounded"
                            />
                        </div>
                        <div>
                            <label htmlFor="dueDate">Due Date:</label>
                            <input
                                type="date"
                                id="dueDate"
                                name="dueDate"
                                className="border px-2 py-1 rounded"
                            />
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
  <div>
    <form
      onSubmit={(e) => {
        e.preventDefault();
        const formData = new FormData(e.currentTarget);

        // Prepare the updated task object with form data
        const updatedTask = {
          ...taskToEdit, // Use taskToEdit to maintain the existing task structure
          name: taskToEdit.name, // formData.get("name"),
          description: taskToEdit.description, // formData.get("description"),
          timeEstimate: taskToEdit.timeEstimate, // formData.get("timeEstimate"),
          dueDate: taskToEdit.dueDate, // formData.get("dueDate"),
        };
        
        // Send the PUT request to the backend to update the task
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
            // Update the tasks in the state with the updated task
            setTasks((prev) =>
              prev.map((t) => (t.id === savedTask.id ? savedTask : t))
            );
            setView("view"); // Switch the view to "view" mode
          })
          .catch((err) => {
            console.error("Error updating task:", err);
          });
      }}
    >
      <div>
        <label htmlFor="name">Name:</label>
        <input
          type="text"
          id="name"
          name="name"
          required
          value={taskToEdit.name}
          onChange={(e) => {
            // @ts-ignore
            setTaskToEdit((prev) => ({
              ...prev,
              name: e.target.value,
            }));
          }}
          className="border px-2 py-1 rounded"
        />
      </div>

      <div>
        <label htmlFor="description">Description:</label>
        <textarea
          id="description"
          name="description"
          required
          value={taskToEdit.description}
          onChange={(e) => {
            // @ts-ignore
            setTaskToEdit((prev) => ({
              ...prev,
              description: e.target.value,
            }));
          }}
          className="border px-2 py-1 rounded"
        />
      </div>

      <div>
        <label htmlFor="timeEstimate">Time Estimate:</label>
        <input
          type="text"
          id="timeEstimate"
          name="timeEstimate"
          value={taskToEdit.timeEstimate}
          onChange={(e) => {
            // @ts-ignore
            setTaskToEdit((prev) => ({
              ...prev,
              timeEstimate: e.target.value,
            }));
          }}
          className="border px-2 py-1 rounded"
        />
      </div>

      <div>
        <label htmlFor="dueDate">Due Date:</label>
        <input
          type="date"
          id="dueDate"
          name="dueDate"
          value={taskToEdit.dueDate}
          onChange={(e) => {
            // @ts-ignore
            setTaskToEdit((prev) => ({
              ...prev,
              dueDate: e.target.value,
            }));
          }}
          className="border px-2 py-1 rounded"
        />
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


        </>
    )
}
