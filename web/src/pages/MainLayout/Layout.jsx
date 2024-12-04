import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import {
  HomeOutlined,
  MacCommandOutlined,
  SignatureOutlined,
  MenuOutlined,
  UserOutlined
} from "@ant-design/icons";
import { Button, Drawer } from "antd";
import logo from "../assets/Logo1.svg";
import backgroundImage from "../assets/Element.svg";
import { useMediaQuery } from "react-responsive";
import { useNavigate } from "react-router-dom";

const Sidebar = ({ children }) => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const isMobileOrTablet = useMediaQuery({ maxWidth: 768 });
  const location = useLocation(); // Untuk mendapatkan URL saat ini

  const navigate = useNavigate()
  // Functions to control drawer
  const showDrawer = () => setDrawerOpen(true);
  const closeDrawer = () => setDrawerOpen(false);

  const handlelogout = () => {
    localStorage.removeItem("authToken")
   navigate("/")
  }

  // Render menu items
  const renderMenu = () => (
    <div className="flex flex-col gap-4">
      <Link
        to="/dashboard"
        className={`flex items-center gap-2 font-Poppins font-medium ${
          location.pathname === "/dashboard" ? "text-blue-500" : "text-black"
        } hover:text-blue-500`}
        onClick={closeDrawer}
      >
        <HomeOutlined /> Dashboard
      </Link>
      <Link
        to="/competence"
        className={`flex items-center gap-2 font-Poppins font-medium ${
          location.pathname === "/competence" ? "text-blue-500" : "text-black"
        } hover:text-blue-500`}
        onClick={closeDrawer}
      >
        <MacCommandOutlined /> Kompetensi
      </Link>
      <Link
        to="/signature"
        className={`flex items-center gap-2 font-Poppins font-medium ${
          location.pathname === "/signature" ? "text-blue-500" : "text-black"
        } hover:text-blue-500`}
        onClick={closeDrawer}
      >
        <SignatureOutlined /> Paraf
      </Link>
    </div>
  );

  return (
    <div
      className="p-9 bg-Background min-h-screen"
      style={{
        backgroundImage: `url(${backgroundImage})`,
        backgroundRepeat: "no-repeat",
        backgroundSize: "cover",
        backgroundPosition: "center",
      }}
    >
      {isMobileOrTablet ? (
        <>
          <Button type="primary" icon={<MenuOutlined />} onClick={showDrawer} />
          <Drawer
            title={<img src={logo} alt="Logo" style={{ width: "100px" }} />}
            placement="left"
            onClose={closeDrawer}
            open={drawerOpen}
            width="50%" // Adjust drawer width as needed
          >
            {renderMenu()}
          </Drawer>
          <div className="mt-4">{children || <p>Main content</p>}</div>
        </>
      ) : (
        <div className="flex align-middle">
          <div className="bg-white flex flex-col w-fit p-8 rounded-xl h-screen font-Poppins text-xl gap-3 justify-between">
           <div className="">
            <div className="my-2">
              <img src={logo} alt="logo" />
            </div>
            {renderMenu()}
           </div>
            <div className="text-red-600 font-Poppins font-medium ">
              <Button
              onClick={() => handlelogout()}
              >
              <UserOutlined /> Logout
              </Button>
            </div>
          </div>
          <div className="flex-1 align-middle justify-center flex items-center">
            {children || <p>Main content</p>}
          </div>
        </div>
      )}

    </div>
  );
};

export default Sidebar;
