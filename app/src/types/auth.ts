export interface User {
  id: string
  email: string
  name: string
  role: 'customer' | 'admin'
}

export interface AuthResponse {
  user: User
}

export interface ApiErrorBody {
  error: {
    code: string
    message: string
  }
}
