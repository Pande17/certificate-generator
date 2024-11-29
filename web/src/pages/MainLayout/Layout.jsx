import { useState } from 'react';
import { Link } from 'react-router-dom';
import { HomeFilled, FileFilled, ToolFilled, MenuOutlined, ToolTwoTone, FileAddTwoTone, HomeTwoTone, SignatureOutlined, SignatureFilled, ToolOutlined, HomeOutlined, FileAddOutlined, MacCommandOutlined } from '@ant-design/icons';
import { Button, Drawer } from 'antd';
import logo from '../assets/Logo1.svg';
import backgroundImage from '../assets/Element.svg';
import { useMediaQuery } from 'react-responsive';

const Sidebar = ({ children }) => {
	const [drawerOpen, setDrawerOpen] = useState(false);
	const isMobileOrTablet = useMediaQuery({ maxWidth: 768 });

	// Functions to control drawer
	const showDrawer = () => setDrawerOpen(true);
	const closeDrawer = () => setDrawerOpen(false);

	// Render menu items
	const renderMenu = () => (
		<div className="flex flex-col gap-4">
			<Link to="/dashboard" className="flex items-center gap-2  font-Poppins font-medium" onClick={closeDrawer}>
				<HomeOutlined /> Dashboard
			</Link>
			<Link to="/competence" className="flex items-center gap-2 font-Poppins font-medium" onClick={closeDrawer}>
				<MacCommandOutlined /> Kompetensi
			</Link>
			<Link to="/signature" className="flex items-center gap-2 font-Poppins font-medium" onClick={closeDrawer}>
				<SignatureOutlined  /> Paraf
			</Link>
		</div>
	);

	return (
		<div
			className="p-9 bg-Background min-h-screen"
			style={{
				backgroundImage: `url(${backgroundImage})`,
				backgroundRepeat: 'no-repeat',
				backgroundCover: 'cover',
				backgroundPosition: 'center',
			}}
		>
			{isMobileOrTablet ? (
				<>
					<Button type="primary" icon={<MenuOutlined />} onClick={showDrawer} />
					<Drawer
						title={<img src={logo} alt="Logo" style={{ width: '100px' }} />}
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
					<div className="flex-1 align-middle justify-center flex items-center">{children || <p>Main content</p>}</div>
				</div>
			)}
		</div>
	);
};

export default Sidebar;
