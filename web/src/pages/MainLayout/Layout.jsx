import { useState } from "react";
import { Link } from "react-router-dom";
import {
  HomeFilled,
  FileFilled,
  ToolFilled,
  MenuOutlined,
  ToolTwoTone,
  FileAddTwoTone,
  HomeTwoTone,
} from "@ant-design/icons";
import { Button, Drawer } from "antd";
import logo from "../assets/logo1.svg";
import backgroundImage from "../assets/Element.svg";
import { useMediaQuery } from "react-responsive";

const Sidebar = ({ children }) => {
  const [drawerOpen, setDrawerOpen] = useState(false);
  const isMobileOrTablet = useMediaQuery({ maxWidth: 768 });

  // Functions to control drawer
  const showDrawer = () => setDrawerOpen(true);
  const closeDrawer = () => setDrawerOpen(false);

  // Render menu items
  const renderMenu = () => (
    <div className="flex flex-col gap-4">
      <Link
        to="/dashboard"
        className="flex items-center gap-2  font-Poppins font-medium"
        onClick={closeDrawer}
      >
        <HomeTwoTone /> Dashboard
      </Link>
      <Link
        to="/create"
        className="flex items-center gap-2  font-Poppins font-medium"
        onClick={closeDrawer}
      >
        <FileAddTwoTone /> Create
      </Link>
      <Link
        to="/tool"
        className="flex items-center gap-2 font-Poppins font-medium"
        onClick={closeDrawer}
      >
        <ToolTwoTone /> Tool
      </Link>
    </div>
  );

  return (
    <div
      className="p-9 bg-Background"
      style={{
        backgroundImage: `url(${backgroundImage})`,
     backgroundRepeat:"no-repeat",
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
          <div className="bg-white flex flex-col w-fit p-8 rounded-xl h-screen font-Poppins text-xl gap-3">
            <div>
              <img src={logo} alt="logo" />
            </div>
            {renderMenu()}
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
