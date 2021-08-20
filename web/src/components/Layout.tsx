import Footer from "./Footer";

type LayoutProps = {
	children: React.ReactNode;
};

const Layout: React.FC<LayoutProps> = ({ children }) => {
	return (
		<div>
			{children}
			<Footer />
		</div>
	);
};

export default Layout;
