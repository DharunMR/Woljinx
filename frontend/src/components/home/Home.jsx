import {useState, useEffect} from 'react';
import axiosClient from '../../api/axiosConfig'
import Movies from '../movies/Movies';
import Spinner from '../spinner/Spinner';

const Home =() => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState();

    useEffect(() => {
        const fetchMovies = async () => {
            setLoading(true);
            setMessage("");
            try{
                console.log("About to call /movies");

                const response = await axiosClient.get('/movies');

                console.log("FULL RESPONSE:", response);
                console.log("RESPONSE.DATA:", response.data);
                console.log("MOVIES ARRAY:", response.data.data);

                const moviesArray = response.data.data;
                setMovies(moviesArray);
                if (response.data.length === 0){
                    setMessage('There are currently no movies available')
                }

            }catch(error){
                console.error('Error fetching movies:', error)
                setMessage("Error fetching movies")
            }finally{
                setLoading(false)
            }
        }
        fetchMovies();
    }, []);

    return (
        <>
            {loading ? (
                <Spinner/>
            ):  (
                <Movies movies ={movies} message ={message}/>
            )}
        </>

    );

};

export default Home;