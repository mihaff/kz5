import pandas as pd
from sklearn.model_selection import train_test_split, GridSearchCV
from sklearn.pipeline import Pipeline
from sklearn.compose import ColumnTransformer
from sklearn.impute import SimpleImputer
from sklearn.preprocessing import OneHotEncoder, StandardScaler
from sklearn.linear_model import LinearRegression
from sklearn.svm import SVR
from sklearn.metrics import mean_squared_error, mean_absolute_error, r2_score
import joblib
import sys


def load_data(file_path):
    if file_path.endswith(".csv"):
        return pd.read_csv(file_path)
    elif file_path.endswith(".xls") or file_path.endswith(".xlsx"):
        return pd.read_excel(file_path)
    elif file_path.endswith(".pkl"):
        return pd.read_pickle(file_path)
    else:
        raise ValueError(
            "Unsupported file type. Supported types are csv, xls, xlsx, and pkl."
        )


def train_model(data, target_column, model_path, algorithm="linear_regression"):
    X = data.drop(columns=[target_column])
    y = data[target_column]

    X_train, X_test, y_train, y_test = train_test_split(
        X, y, test_size=0.2, random_state=42
    )

    categorical_features = X.select_dtypes(include=["object"]).columns
    numeric_features = X.select_dtypes(include=["int", "float"]).columns

    categorical_transformer = Pipeline(
        steps=[
            ("imputer", SimpleImputer(strategy="most_frequent")),
            ("onehot", OneHotEncoder(handle_unknown="ignore")),
        ]
    )

    numeric_transformer = Pipeline(
        steps=[
            ("imputer", SimpleImputer(strategy="mean")),
            ("scaler", StandardScaler()),
        ]
    )

    preprocessor = ColumnTransformer(
        transformers=[
            ("num", numeric_transformer, numeric_features),
            ("cat", categorical_transformer, categorical_features),
        ]
    )

    if algorithm == "linear_regression":
        regressor = LinearRegression()
        param_grid = {}  # "regressor__normalize": [True, False]}
    elif algorithm == "support_vector_machine":
        regressor = SVR()
        param_grid = {
            # "regressor__kernel": ["linear", "poly", "rbf", "sigmoid"],
            # "regressor__C": [0.01, 0.1, 1, 10, 100],
            # "regressor__gamma": ["scale", "auto", 0.001, 0.01, 0.1, 1],
            # "regressor__epsilon": [0.01, 0.1, 1],
        }
    else:
        raise ValueError(
            "Unsupported algorithm. Choose 'linear_regression' or 'support_vector_machine'."
        )

    clf = Pipeline(steps=[("preprocessor", preprocessor), ("regressor", regressor)])

    grid_search = GridSearchCV(clf, param_grid, cv=5, scoring="neg_mean_squared_error")
    grid_search.fit(X_train, y_train)

    y_pred = grid_search.predict(X_test)
    mse = mean_squared_error(y_test, y_pred)
    rmse = mse**0.5
    mae = mean_absolute_error(y_test, y_pred)
    r2 = r2_score(y_test, y_pred)

    print(f"RMSE: {rmse:.2f}")
    print(f"MAE: {mae:.2f}")
    print(f"R2 Score: {r2:.2f}")

    joblib.dump(grid_search.best_estimator_, model_path)

    return grid_search.best_estimator_


def main():
    if len(sys.argv) != 5:
        print(
            "Usage: python model_class.py <target_column> <input_file_path> <output_file_path>"
        )
        sys.exit(1)

    algorithm = sys.argv[1]
    target_column = sys.argv[2]
    input_file_path = sys.argv[3]
    output_file_path = sys.argv[4]

    data = load_data(input_file_path)
    train_model(data, target_column, algorithm=algorithm, model_path=output_file_path)


if __name__ == "__main__":
    main()
