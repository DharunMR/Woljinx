import Movie from "../movie/Movie";

const Movies = ({ movies, updateMovieReview, message }) => {
    const displayMessage = message || "No movies available";

    return (
        <div className="container mt-4">
            <div className="row">

                {Array.isArray(movies) && movies.length > 0 ? (
                    movies
                        .filter(Boolean)
                        .map(movie => (
                            <Movie
                                key={movie._id || movie.imdb_id}
                                movie={movie}
                                updateMovieReview={updateMovieReview}
                            />
                        ))
                ) : (
                    <h2 className="text-center">{displayMessage}</h2>
                )}

            </div>
        </div>
    );
};

export default Movies;
