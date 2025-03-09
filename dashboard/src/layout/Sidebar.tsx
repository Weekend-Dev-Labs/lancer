import { useNavigate } from "react-router-dom";

const Sidebar = () => {
  const navigate = useNavigate();
  const menuItems = [
    { title: "Dashboard", path: "" },
    { title: "Uploads", path: "/upload" },
  ];

  return (
    <div className="h-full w-64 bg-white shadow-md flex flex-col p-5">
      <h2 className="text-2xl font-bold mb-4">Admin Panel</h2>
      <nav className="flex flex-col gap-3">
        {menuItems.map((item, index) => (
          <button
            onClick={() => navigate(item.path)}
            key={index}
            className="text-lg text-gray-700 hover:bg-gray-100 px-4 py-2 rounded-lg text-left"
          >
            {item.title}
          </button>
        ))}
      </nav>
    </div>
  );
};

export default Sidebar;
