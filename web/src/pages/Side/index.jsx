import { useState } from 'react';
import { Link } from 'react-router-dom';
import { HomeFilled, FileFilled, ToolFilled, MenuOutlined } from '@ant-design/icons';
import { Button, Drawer } from 'antd';
import logo from '../assets/Logo1.svg';
import { useMediaQuery } from 'react-responsive';

const Sidebar = () => {
	const [drawerOpen, setDrawerOpen] = useState(false);
	const isMobileOrTablet = useMediaQuery({ maxWidth: 768 });

	// Functions to control drawer
	const showDrawer = () => setDrawerOpen(true);
	const closeDrawer = () => setDrawerOpen(false);

	// Render menu items
	const renderMenu = () => (
		<div className="flex flex-col gap-4">
			<Link to="/dashboard" className="flex items-center gap-2" onClick={closeDrawer}>
				<HomeFilled /> Dashboard
			</Link>
			<Link to="/create" className="flex items-center gap-2" onClick={closeDrawer}>
				<FileFilled /> Create
			</Link>
			<Link to="/tool" className="flex items-center gap-2" onClick={closeDrawer}>
				<ToolFilled /> Tool
			</Link>
		</div>
	);

	return (
		<>
			{isMobileOrTablet ? (
				<>
					<Button type="primary" icon={<MenuOutlined />} onClick={showDrawer} />
					<Drawer title={<img src={logo} alt="Logo" style={{ width: '100px' }} />} placement="left" onClose={closeDrawer} open={drawerOpen}>
						{renderMenu()}
					</Drawer>
				</>
			) : (
				<div className="p-9 bg-Background">
					<div className="flex align-middle ">
						<div className="bg-white flex flex-col w-fit p-8 rounded-xl h-screen font-Poppins text-xl gap-3 ">
							<div>
								<img src={logo} alt="logo" />
							</div>
							{renderMenu()}
						</div>
						<div className="w-screen align-middle justify-center bg-slate-400 flex items-center">
							<p>main content</p>
						</div>
					</div>
				</div>
			)}
		</>
	);
};

export default Sidebar;
