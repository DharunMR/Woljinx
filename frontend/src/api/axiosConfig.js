import axios from 'axios';
// const apiUrl = import.meta.env.VITE_API_BASE_URL;

export default axios.create({
    baseURL:"http://backendapp-service:8080",
    headers:{'Content-Type':'application/json'},
    withCredentials: true,
})