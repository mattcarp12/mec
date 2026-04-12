import { BrowserRouter, Routes, Route, Link, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { api } from './api';
import { AuthProvider, useAuth } from './AuthContext';

// --- HOME COMPONENT ---
function Home() {
  const { isAuthenticated, logout } = useAuth();

  return (
    <div>
      <h2>Product Catalog</h2>
      <p>Welcome to our store! Browse our items below.</p>

      {isAuthenticated ? (
        <div style={{ padding: '10px', backgroundColor: '#eef' }}>
          <strong>You are logged in!</strong> The secure cookie is stored in your browser.
          <br /><button onClick={logout}>Log Out (Local only)</button>
        </div>
      ) : (
        <p><em>You are currently browsing as a guest.</em></p>
      )}
    </div>
  );
}

// --- LOGIN COMPONENT ---
function Login() {
  const { register, handleSubmit, formState: { errors } } = useForm();
  const { login } = useAuth();
  const navigate = useNavigate();

  // This function runs when the form is submitted
  const onSubmit = async (data: any) => {
    try {
      // 1. Send the data to our Go Traefik Gateway
      await api.post('/api/v1/auth/login', data);

      // 2. If successful, update our global React state
      login();

      // 3. Redirect the user back to the catalog
      navigate('/');
    } catch (error) {
      alert("Login failed. Check console or ensure Go backend is running.");
      console.error(error);
    }
  };

  return (
    <div>
      <h2>Login</h2>
      {/* handleSubmit prevents the default HTML page refresh */}
      <form onSubmit={handleSubmit(onSubmit)}>
        <div>
          <label>Email: </label>
          {/* ...register connects this input to react-hook-form */}
          <input type="email" {...register("email", { required: true })} />
          {errors.email && <span style={{ color: 'red' }}> Required</span>}
        </div>

        <div style={{ marginTop: '10px' }}>
          <label>Password: </label>
          <input type="password" {...register("password", { required: true })} />
          {errors.password && <span style={{ color: 'red' }}> Required</span>}
        </div>

        <button type="submit" style={{ marginTop: '15px' }}>Submit</button>
      </form>
    </div>
  );
}

// --- REGISTER COMPONENT (Placeholder) ---
function Register() {
  return <div><h2>Register</h2><p>Registration form goes here.</p></div>;
}

// 3. This is the main App component that we imported into main.tsx.
export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>

        {/* This is our persistent navigation bar. It will ALWAYS show at the top 
          because it is outside of the <Routes> block. */}
        <nav>
          <h1>Matt's E-Commerce</h1>
          {/* <Link> is React's version of an <a> tag. It changes the URL 
            without reloading the browser. */}
          <Link to="/">Home</Link> |{' '}
          <Link to="/login">Login</Link> |{' '}
          <Link to="/register">Register</Link>
        </nav>

        <hr />

        {/* The Routes block acts like a switchboard. It looks at the URL and 
          renders ONLY the Route that matches. */}
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
        </Routes>

      </BrowserRouter>
    </AuthProvider>
  )
}