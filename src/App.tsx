import { BrowserRouter as Router } from "react-router-dom";
import AppRouter from "routes";
import "./App.css";

const App: React.FC = () => {
	return (
		<Router>
			<div className="App">
				<AppRouter />
			</div>
		</Router>
	);
};

export default App;
