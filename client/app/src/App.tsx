import './App.css'
import {routes} from "./routers/route.tsx";
import {RouterProvider} from "react-router-dom";

function App() {
  return (
    <RouterProvider router={routes}/>
  )
}

export default App
