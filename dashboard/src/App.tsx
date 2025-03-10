import Login from "./auth/Login";
import Layout from "./layout/Layout";
import QueryProvider from "./lib/tanStack";
import { Route, Routes } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import UploadList from "./pages/UploadList";
import Session from "./pages/Session";
import Settings from "./pages/Settings";

const App = () => {
  return (
    <QueryProvider>
      <Routes>
        <Route path="login" element={<Login />} />
        <Route path="/" element={<Layout />}>
          <Route index element={<Dashboard />} />
          <Route path="/upload" element={<UploadList />} />
          <Route path="/session" element={<Session />} />
          <Route path="/settings" element={<Settings />} />
        </Route>
      </Routes>
    </QueryProvider>
  );
};

export default App;
