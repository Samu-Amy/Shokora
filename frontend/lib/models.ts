type LoginUserRes = {
  user: User;
  verification_id: string;
  is_email_sent: boolean;
};

type User = {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  birthday: Date;
  is_verified: boolean;
  role: UserRole;
  permissions: number;
  created_at: Date;
  updated_at: Date;
};
