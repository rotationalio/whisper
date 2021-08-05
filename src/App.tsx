import Footer from "components/Footer";
import { BrowserRouter as Router } from "react-router-dom";
import AppRouter from "routes";
import "./App.css";

const App: React.FC = () => {
	return (
		<Router>
			<div className="App">
				<AppRouter />
				<Footer />
			</div>
		</Router>
	);
};

export default App;
