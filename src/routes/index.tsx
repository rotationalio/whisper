import NotFound from "pages/NotFound";
import React, { Suspense } from "react";
import { Route, Switch } from "react-router-dom";

const ShowSecret = React.lazy(() => import("pages/FetchSecret"));
const CreateSecret = React.lazy(() => import("pages/CreateSecret"));

const AppRouter: React.FC = () => {
	return (
		<Suspense fallback="loading...">
			<Switch>
				<Route path="/secret/:token" exact component={ShowSecret} />
				<Route path="/" exact component={CreateSecret} />
				<Route path="*" component={NotFound} />
			</Switch>
		</Suspense>
	);
};

export default AppRouter;
